// Copyright 2020-2026 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buflsp

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf16"

	"github.com/bufbuild/protocompile/experimental/protoscope"
	"go.lsp.dev/protocol"
)

// protoscopeManager tracks open .protoscope files in the LSP session.
type protoscopeManager struct {
	lsp *lsp

	mu        sync.Mutex
	uriToFile map[protocol.URI]*protoscopeFile
}

type protoscopeFile struct {
	text    string
	version int32
}

func newProtoscopeManager(lsp *lsp) *protoscopeManager {
	return &protoscopeManager{
		lsp:       lsp,
		uriToFile: make(map[protocol.URI]*protoscopeFile),
	}
}

// Track opens or updates a .protoscope file.
func (m *protoscopeManager) Track(ctx context.Context, uri protocol.URI, version int32, text string) {
	normalized := normalizeURI(uri)
	m.mu.Lock()
	m.uriToFile[normalized] = &protoscopeFile{
		text:    text,
		version: version,
	}
	m.mu.Unlock()

	m.checkAndPublishDiagnostics(ctx, normalized, version, text)
}

// Close closes a .protoscope file.
func (m *protoscopeManager) Close(ctx context.Context, uri protocol.URI) {
	normalized := normalizeURI(uri)
	m.mu.Lock()
	delete(m.uriToFile, normalized)
	m.mu.Unlock()

	// Clear diagnostics on client
	_ = m.lsp.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         normalized,
		Diagnostics: []protocol.Diagnostic{},
	})
}

// Get retrieves a .protoscope file if it exists.
func (m *protoscopeManager) Get(uri protocol.URI) (*protoscopeFile, bool) {
	normalized := normalizeURI(uri)
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.uriToFile[normalized]
	return f, ok
}

// Has reports whether the URI is tracked as a protoscope file.
func (m *protoscopeManager) Has(uri protocol.URI) bool {
	_, ok := m.Get(uri)
	return ok
}

// checkAndPublishDiagnostics runs diagnostics on the protoscope text and publishes them.
func (m *protoscopeManager) checkAndPublishDiagnostics(ctx context.Context, uri protocol.URI, version int32, text string) {
	diags := protoscope.Diagnostics(uri.Filename(), []byte(text))
	protocolDiags := make([]protocol.Diagnostic, len(diags))
	for i, d := range diags {
		var severity protocol.DiagnosticSeverity
		switch d.Level {
		case protoscope.SeverityInfo:
			severity = protocol.DiagnosticSeverityInformation
		case protoscope.SeverityWarning:
			severity = protocol.DiagnosticSeverityWarning
		case protoscope.SeverityError:
			severity = protocol.DiagnosticSeverityError
		default:
			severity = protocol.DiagnosticSeverityError
		}

		protocolDiags[i] = protocol.Diagnostic{
			Source:   serverName,
			Severity: severity,
			Message:  d.Message,
			Range:    protoscopeRangeToProtocolRange(d.Range),
		}
	}

	_ = m.lsp.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Version:     uint32(version),
		Diagnostics: protocolDiags,
	})
}

// GetHover retrieves hover documentation for a position in a protoscope file.
func (m *protoscopeManager) GetHover(ctx context.Context, uri protocol.URI, pos protocol.Position) (*protocol.Hover, error) {
	file, ok := m.Get(uri)
	if !ok {
		return nil, nil
	}

	// protoscope.Hover expects 1-indexed line and column
	line := int(pos.Line) + 1
	col := int(pos.Character) + 1

	h, err := protoscope.Hover(uri.Filename(), []byte(file.text), line, col)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, nil
	}

	rangeVal := protoscopeRangeToProtocolRange(h.Range)
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: h.Text,
		},
		Range: &rangeVal,
	}, nil
}

// GetDocumentSymbols retrieves the symbols within a protoscope file.
func (m *protoscopeManager) GetDocumentSymbols(ctx context.Context, uri protocol.URI) ([]any, error) {
	file, ok := m.Get(uri)
	if !ok {
		return nil, nil
	}

	symbols, _ := protoscope.DocumentSymbols(uri.Filename(), []byte(file.text))
	docSymbols := convertProtoscopeDocumentSymbols(symbols)

	anyResults := make([]any, len(docSymbols))
	for i, ds := range docSymbols {
		anyResults[i] = ds
	}
	return anyResults, nil
}

// Formatting formats a protoscope file.
func (m *protoscopeManager) Formatting(ctx context.Context, uri protocol.URI) ([]protocol.TextEdit, error) {
	file, ok := m.Get(uri)
	if !ok {
		return nil, nil
	}

	binary, diags := protoscope.Assemble(uri.Filename(), []byte(file.text))
	var hasError bool
	for _, d := range diags {
		if d.Level == protoscope.SeverityError {
			hasError = true
			break
		}
	}
	if hasError {
		return nil, nil // Do not format if there are compiler errors
	}

	newText, err := protoscope.Disassemble(binary, protoscope.DisassembleOptions{})
	if err != nil {
		return nil, err
	}

	if newText == file.text {
		return nil, nil
	}

	// Calculate end position
	endLine := strings.Count(file.text, "\n")
	endCharacter := 0
	lastNewLine := strings.LastIndexByte(file.text, '\n')
	var lastLine string
	if lastNewLine == -1 {
		lastLine = file.text
	} else {
		lastLine = file.text[lastNewLine+1:]
	}
	for _, r := range lastLine {
		endCharacter += utf16.RuneLen(r)
	}

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 0,
				},
				End: protocol.Position{
					Line:      uint32(endLine),
					Character: uint32(endCharacter),
				},
			},
			NewText: newText,
		},
	}, nil
}

// protoscopeRangeToProtocolRange converts a 1-indexed protoscope.Range to a 0-indexed protocol.Range.
func protoscopeRangeToProtocolRange(r protoscope.Range) protocol.Range {
	startLine := uint32(0)
	if r.Start.Line > 0 {
		startLine = uint32(r.Start.Line - 1)
	}
	startCol := uint32(0)
	if r.Start.Column > 0 {
		startCol = uint32(r.Start.Column - 1)
	}
	endLine := uint32(0)
	if r.End.Line > 0 {
		endLine = uint32(r.End.Line - 1)
	}
	endCol := uint32(0)
	if r.End.Column > 0 {
		endCol = uint32(r.End.Column - 1)
	}

	return protocol.Range{
		Start: protocol.Position{
			Line:      startLine,
			Character: startCol,
		},
		End: protocol.Position{
			Line:      endLine,
			Character: endCol,
		},
	}
}

// convertProtoscopeDocumentSymbols recursively converts protoscope.DocumentSymbols to protocol.DocumentSymbols.
func convertProtoscopeDocumentSymbols(symbols []protoscope.DocumentSymbol) []protocol.DocumentSymbol {
	if len(symbols) == 0 {
		return nil
	}

	res := make([]protocol.DocumentSymbol, len(symbols))
	for i, s := range symbols {
		var kind protocol.SymbolKind
		switch s.Kind {
		case "field":
			kind = protocol.SymbolKindField
		case "block":
			kind = protocol.SymbolKindNamespace
		case "literal":
			kind = protocol.SymbolKindConstant
		default:
			kind = protocol.SymbolKindField
		}

		r := protoscopeRangeToProtocolRange(s.Range)

		res[i] = protocol.DocumentSymbol{
			Name:           s.Name,
			Detail:         s.Detail,
			Kind:           kind,
			Range:          r,
			SelectionRange: r,
			Children:       convertProtoscopeDocumentSymbols(s.Children),
		}
	}
	return res
}

// commandDisassembleProtoscope is the command to disassemble a binary file to protoscope text format.
const commandDisassembleProtoscope = "buf.protoscope.disassemble.server"

// ExecuteDisassemble reads the binary file at the given URI and disassembles it to protoscope text format.
func (m *protoscopeManager) ExecuteDisassemble(ctx context.Context, uri protocol.URI) (string, error) {
	filename := uri.Filename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return protoscope.Disassemble(data, protoscope.DisassembleOptions{})
}

// commandAssembleProtoscope is the command to assemble protoscope text format to binary.
const commandAssembleProtoscope = "buf.protoscope.assemble.server"

// ExecuteAssemble reads the protoscope file (either from the tracker or from disk)
// and compiles it to a binary protobuf file. Returns the binary content encoded as a base64 string.
func (m *protoscopeManager) ExecuteAssemble(ctx context.Context, uri protocol.URI) (string, error) {
	var text string
	if file, ok := m.Get(uri); ok {
		text = file.text
	} else {
		filename := uri.Filename()
		data, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}
		text = string(data)
	}

	binary, diags := protoscope.Assemble(uri.Filename(), []byte(text))
	var hasError bool
	var errorMessages []string
	for _, d := range diags {
		if d.Level == protoscope.SeverityError {
			hasError = true
			errorMessages = append(errorMessages, d.Message)
		}
	}
	if hasError {
		return "", fmt.Errorf("assembly failed with errors:\n%s", strings.Join(errorMessages, "\n"))
	}

	return base64.StdEncoding.EncodeToString(binary), nil
}
