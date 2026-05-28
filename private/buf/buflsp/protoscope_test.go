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

package buflsp_test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.lsp.dev/protocol"
)

func TestProtoscopeDiagnostics(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "valid.protoscope")
		content := `1: 150
2: {
  3: "hello"
}
`
		err := os.WriteFile(filePath, []byte(content), 0600)
		require.NoError(t, err)

		_, testURI, capture := setupLSPServerWithDiagnostics(t, filePath)

		diagnostics := capture.wait(t, testURI, 5*time.Second, func(p *protocol.PublishDiagnosticsParams) bool {
			return true
		})
		require.NotNil(t, diagnostics)
		assert.Empty(t, diagnostics.Diagnostics)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.protoscope")
		content := `1: 
2: {
`
		err := os.WriteFile(filePath, []byte(content), 0600)
		require.NoError(t, err)

		_, testURI, capture := setupLSPServerWithDiagnostics(t, filePath)

		diagnostics := capture.wait(t, testURI, 5*time.Second, func(p *protocol.PublishDiagnosticsParams) bool {
			return len(p.Diagnostics) > 0
		})
		require.NotNil(t, diagnostics)
		assert.NotEmpty(t, diagnostics.Diagnostics)
		assert.Equal(t, protocol.DiagnosticSeverityError, diagnostics.Diagnostics[0].Severity)
	})
}

func TestProtoscopeHover(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.protoscope")
	content := `1: 150
2: {
  3: ` + "`" + `01 02 03` + "`" + `
}
`
	err := os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	clientJSONConn, testURI := setupLSPServer(t, filePath)

	// Hover over field tag "1:" (line 0, character 0)
	var hover *protocol.Hover
	_, hoverErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentHover, protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Position: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
	}, &hover)
	require.NoError(t, hoverErr)
	require.NotNil(t, hover)
	assert.Contains(t, hover.Contents.Value, "Field Number")
	assert.Contains(t, hover.Contents.Value, "1")

	// Hover over literal value "150" (line 0, character 3)
	var hover2 *protocol.Hover
	_, hoverErr2 := clientJSONConn.Call(ctx, protocol.MethodTextDocumentHover, protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Position: protocol.Position{
				Line:      0,
				Character: 3,
			},
		},
	}, &hover2)
	require.NoError(t, hoverErr2)
	require.NotNil(t, hover2)
	assert.Contains(t, hover2.Contents.Value, "Literal Value")
	assert.Contains(t, hover2.Contents.Value, "150")
}

func TestProtoscopeDocumentSymbol(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "symbols.protoscope")
	content := `1: 150
2: {
  3: "hello"
}
`
	err := os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	clientJSONConn, testURI := setupLSPServer(t, filePath)

	var symbols []protocol.DocumentSymbol
	_, symErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentDocumentSymbol, protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: testURI,
		},
	}, &symbols)
	require.NoError(t, symErr)
	require.Len(t, symbols, 2)

	// First symbol: Field 1
	assert.Equal(t, "1:", symbols[0].Name)
	assert.Equal(t, protocol.SymbolKindField, symbols[0].Kind)

	// Second symbol: Field 2
	assert.Equal(t, "2:", symbols[1].Name)
	assert.Equal(t, protocol.SymbolKindField, symbols[1].Kind)
	require.Len(t, symbols[1].Children, 1)

	// Children of Field 2: block "Length-Prefixed"
	blockSymbol := symbols[1].Children[0]
	assert.Equal(t, "Length-Prefixed", blockSymbol.Name)
	assert.Equal(t, protocol.SymbolKindNamespace, blockSymbol.Kind)
	require.Len(t, blockSymbol.Children, 1)

	// Inside Length-Prefixed: Field 3
	field3 := blockSymbol.Children[0]
	assert.Equal(t, "3:", field3.Name)
	assert.Equal(t, protocol.SymbolKindField, field3.Kind)
}

func TestProtoscopeFormatting(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	t.Run("unformatted", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "unformatted.protoscope")
		content := `1:150
2:{
  3:"hello"
}`
		err := os.WriteFile(filePath, []byte(content), 0600)
		require.NoError(t, err)

		clientJSONConn, testURI := setupLSPServer(t, filePath)

		var textEdits []protocol.TextEdit
		_, formatErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentFormatting, protocol.DocumentFormattingParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
		}, &textEdits)
		require.NoError(t, formatErr)
		require.Len(t, textEdits, 1)

		assert.Equal(t, uint32(0), textEdits[0].Range.Start.Line)
		assert.Equal(t, uint32(0), textEdits[0].Range.Start.Character)
		assert.Contains(t, textEdits[0].NewText, "1: 150")
	})

	t.Run("invalid_no_format", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.protoscope")
		content := `1: 
2: {
`
		err := os.WriteFile(filePath, []byte(content), 0600)
		require.NoError(t, err)

		clientJSONConn, testURI := setupLSPServer(t, filePath)

		var textEdits []protocol.TextEdit
		_, formatErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentFormatting, protocol.DocumentFormattingParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
		}, &textEdits)
		require.NoError(t, formatErr)
		assert.Empty(t, textEdits)
	})
}

func TestProtoscopeDisassembleCommand(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.bin")
	// Protobuf binary with field 1 = 150 (varint)
	binaryData := []byte{0x08, 0x96, 0x01}
	err := os.WriteFile(filePath, binaryData, 0600)
	require.NoError(t, err)

	clientJSONConn, testURI := setupLSPServer(t, filePath)

	var disassembledText string
	_, cmdErr := clientJSONConn.Call(ctx, protocol.MethodWorkspaceExecuteCommand, protocol.ExecuteCommandParams{
		Command:   "buf.protoscope.disassemble.server",
		Arguments: []any{string(testURI)},
	}, &disassembledText)
	require.NoError(t, cmdErr)
	assert.Contains(t, disassembledText, "1: 150")
}

func TestProtoscopeAssembleCommand(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.protoscope")
	content := "1: 150\n"
	err := os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	clientJSONConn, testURI := setupLSPServer(t, filePath)

	var assembledBase64 string
	_, cmdErr := clientJSONConn.Call(ctx, protocol.MethodWorkspaceExecuteCommand, protocol.ExecuteCommandParams{
		Command:   "buf.protoscope.assemble.server",
		Arguments: []any{string(testURI)},
	}, &assembledBase64)
	require.NoError(t, cmdErr)

	binaryData, err := base64.StdEncoding.DecodeString(assembledBase64)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x08, 0x96, 0x01}, binaryData)
}

func TestProtoscopeAssembleCommandError(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.protoscope")
	content := `1:
2: {
`
	err := os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	clientJSONConn, testURI := setupLSPServer(t, filePath)

	var assembledBase64 string
	_, cmdErr := clientJSONConn.Call(ctx, protocol.MethodWorkspaceExecuteCommand, protocol.ExecuteCommandParams{
		Command:   "buf.protoscope.assemble.server",
		Arguments: []any{string(testURI)},
	}, &assembledBase64)
	assert.Error(t, cmdErr)
}

func TestProtoscopeUntitled(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tempDir := t.TempDir()
	dummyPath := filepath.Join(tempDir, "dummy.proto")
	err := os.WriteFile(dummyPath, []byte("syntax = \"proto3\";"), 0600)
	require.NoError(t, err)

	clientJSONConn, _ := setupLSPServer(t, dummyPath)

	untitledURI := protocol.URI("untitled:/Untitled-1")
	content := `1: 150
2: {
  3: "hello"
}
`

	// 1. Send DidOpen for untitled document with language ID protoscope
	err = clientJSONConn.Notify(ctx, protocol.MethodTextDocumentDidOpen, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        untitledURI,
			LanguageID: "protoscope",
			Version:    1,
			Text:       content,
		},
	})
	require.NoError(t, err)

	// 2. Test Hover on untitled URI
	var hover *protocol.Hover
	_, hoverErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentHover, protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: untitledURI,
			},
			Position: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
	}, &hover)
	require.NoError(t, hoverErr)
	require.NotNil(t, hover)
	assert.Contains(t, hover.Contents.Value, "Field Number")

	// 3. Test Document Symbols on untitled URI
	var symbols []protocol.DocumentSymbol
	_, symErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentDocumentSymbol, protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: untitledURI,
		},
	}, &symbols)
	require.NoError(t, symErr)
	require.Len(t, symbols, 2)

	// 4. Test Formatting on untitled URI
	var textEdits []protocol.TextEdit
	_, formatErr := clientJSONConn.Call(ctx, protocol.MethodTextDocumentFormatting, protocol.DocumentFormattingParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: untitledURI,
		},
	}, &textEdits)
	require.NoError(t, formatErr)

	// 5. Test Execute Command (Assemble) on untitled URI with mismatched slash
	var assembledBase64 string
	_, cmdErr := clientJSONConn.Call(ctx, protocol.MethodWorkspaceExecuteCommand, protocol.ExecuteCommandParams{
		Command:   "buf.protoscope.assemble.server",
		Arguments: []any{"untitled:Untitled-1"},
	}, &assembledBase64)
	require.NoError(t, cmdErr)

	binaryData, err := base64.StdEncoding.DecodeString(assembledBase64)
	require.NoError(t, err)
	assert.NotEmpty(t, binaryData)
}
