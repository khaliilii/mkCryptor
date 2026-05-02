
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/encrypt", encryptHandler)
	http.HandleFunc("/decrypt", decryptHandler)
	
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("🔐 CRYPTO APP - Complete Encryption/Decryption Tool")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("🌐 Server started at: http://localhost:8080")
	fmt.Println("📱 Open this URL in your browser")
	fmt.Println("⚡ Press Ctrl+C to stop the server")
	fmt.Println("═══════════════════════════════════════════════════")
	
	http.ListenAndServe(":8080", nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🔐 CryptoApp - Complete Encryption Tool</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', sans-serif;
            background: linear-gradient(135deg, #0f2027 0%, #203a43 50%, #2c5364 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 20px;
            padding: 30px;
            margin-bottom: 30px;
            text-align: center;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            color: white;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .header p {
            font-size: 1.1em;
            opacity: 0.95;
        }
        
        .row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
            margin-bottom: 30px;
        }
        
        .card {
            background: white;
            border-radius: 20px;
            padding: 25px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            transition: transform 0.3s;
        }
        
        .card:hover {
            transform: translateY(-5px);
        }
        
        .card h2 {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 20px;
            font-size: 1.8em;
        }
        
        .card.encrypt h2 {
            color: #28a745;
        }
        
        .card.decrypt h2 {
            color: #dc3545;
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            color: #555;
            font-weight: 600;
            font-size: 0.9em;
        }
        
        textarea, input {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 14px;
            transition: all 0.3s;
            font-family: 'Courier New', monospace;
        }
        
        textarea:focus, input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        textarea {
            min-height: 120px;
            resize: vertical;
        }
        
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }
        
        button {
            flex: 1;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 20px;
            border: none;
            border-radius: 10px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.2);
        }
        
        button.encrypt-btn {
            background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
        }
        
        button.decrypt-btn {
            background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
        }
        
        button.clear-btn {
            background: linear-gradient(135deg, #4b6cb7 0%, #182848 100%);
        }
        
        .result-area {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 15px;
            margin-top: 15px;
            border-left: 4px solid #667eea;
        }
        
        .result-label {
            font-weight: 600;
            color: #667eea;
            margin-bottom: 8px;
            font-size: 0.85em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .result-content {
            word-break: break-all;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            color: #333;
            margin-bottom: 10px;
            max-height: 150px;
            overflow-y: auto;
        }
        
        .copy-btn {
            background: #28a745;
            padding: 6px 12px;
            font-size: 12px;
            margin-top: 5px;
        }
        
        .status {
            padding: 10px;
            margin-top: 10px;
            border-radius: 8px;
            display: none;
            font-size: 14px;
        }
        
        .status.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
            display: block;
        }
        
        .status.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
            display: block;
        }
        
        .info-card {
            background: white;
            border-radius: 20px;
            padding: 25px;
            margin-top: 0;
        }
        
        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        
        .info-item {
            text-align: center;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
        }
        
        .info-item h4 {
            color: #667eea;
            margin-bottom: 8px;
        }
        
        @media (max-width: 768px) {
            .row {
                grid-template-columns: 1fr;
            }
            
            .header h1 {
                font-size: 1.5em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 CryptoApp Complete</h1>
            <p>AES-256-GCM Military Grade Encryption | 100% Offline | Your Data Never Leaves Your Device</p>
        </div>
        
        <div class="row">
            <!-- Encryption Card -->
            <div class="card encrypt">
                <h2>
                    <span>🔒</span>
                    <span>Encryption</span>
                </h2>
                <div class="form-group">
                    <label>📝 Plain Text</label>
                    <textarea id="encryptText" placeholder="Enter your secret message here..."></textarea>
                </div>
                <div class="form-group">
                    <label>🔑 Encryption Key</label>
                    <input type="password" id="encryptKey" placeholder="Enter encryption key (min 6 characters)">
                </div>
                <div class="button-group">
                    <button onclick="encrypt()" class="encrypt-btn">🔒 Encrypt</button>
                    <button onclick="clearEncrypt()" class="clear-btn">🗑️ Clear</button>
                </div>
                <div id="encryptResultArea" class="result-area" style="display:none;">
                    <div class="result-label">🔐 ENCRYPTED RESULT</div>
                    <div id="encryptResult" class="result-content"></div>
                    <button onclick="copyEncryptResult()" class="copy-btn">📋 Copy to Clipboard</button>
                </div>
                <div id="encryptStatus" class="status"></div>
            </div>
            
            <!-- Decryption Card -->
            <div class="card decrypt">
                <h2>
                    <span>🔓</span>
                    <span>Decryption</span>
                </h2>
                <div class="form-group">
                    <label>🔐 Encrypted Text</label>
                    <textarea id="decryptText" placeholder="Paste encrypted text here..."></textarea>
                </div>
                <div class="form-group">
                    <label>🔑 Decryption Key</label>
                    <input type="password" id="decryptKey" placeholder="Enter the same key used for encryption">
                </div>
                <div class="button-group">
                    <button onclick="decrypt()" class="decrypt-btn">🔓 Decrypt</button>
                    <button onclick="clearDecrypt()" class="clear-btn">🗑️ Clear All</button>
                </div>
                <div id="decryptResultArea" class="result-area" style="display:none;">
                    <div class="result-label">📄 DECRYPTED RESULT</div>
                    <div id="decryptResult" class="result-content"></div>
                    <button onclick="copyDecryptResult()" class="copy-btn">📋 Copy to Clipboard</button>
                </div>
                <div id="decryptStatus" class="status"></div>
            </div>
        </div>
        
        <div class="info-card">
            <h3 style="color: #667eea; margin-bottom: 15px;">ℹ️ About This Tool</h3>
            <div class="info-grid">
                <div class="info-item">
                    <h4>🛡️ Security</h4>
                    <p style="font-size: 13px;">AES-256-GCM + SHA-256</p>
                </div>
                <div class="info-item">
                    <h4>💻 Offline</h4>
                    <p style="font-size: 13px;">100% Local Processing</p>
                </div>
                <div class="info-item">
                    <h4>🔑 Key Required</h4>
                    <p style="font-size: 13px;">Same key for encrypt/decrypt</p>
                </div>
                <div class="info-item">
                    <h4>⚠️ Warning</h4>
                    <p style="font-size: 13px;">Lost key = Lost data</p>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        let currentEncryptResult = '';
        let currentDecryptResult = '';
        
        async function encrypt() {
            const text = document.getElementById('encryptText').value;
            const key = document.getElementById('encryptKey').value;
            
            if (!text || !key) {
                showStatus('encryptStatus', '❌ Please enter both text and key', 'error');
                return;
            }
            
            if (key.length < 6) {
                showStatus('encryptStatus', '❌ Key must be at least 6 characters', 'error');
                return;
            }
            
            try {
                const response = await fetch('/encrypt', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({text, key})
                });
                const data = await response.json();
                
                if (data.error) {
                    showStatus('encryptStatus', '❌ ' + data.error, 'error');
                } else {
                    currentEncryptResult = data.result;
                    document.getElementById('encryptResult').innerHTML = data.result;
                    document.getElementById('encryptResultArea').style.display = 'block';
                    showStatus('encryptStatus', '✅ Text encrypted successfully!', 'success');
                }
            } catch (error) {
                showStatus('encryptStatus', '❌ Error: ' + error.message, 'error');
            }
        }
        
        async function decrypt() {
            const text = document.getElementById('decryptText').value;
            const key = document.getElementById('decryptKey').value;
            
            if (!text || !key) {
                showStatus('decryptStatus', '❌ Please enter both encrypted text and key', 'error');
                return;
            }
            
            try {
                const response = await fetch('/decrypt', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({text, key})
                });
                const data = await response.json();
                
                if (data.error) {
                    showStatus('decryptStatus', '❌ ' + data.error, 'error');
                } else {
                    currentDecryptResult = data.result;
                    document.getElementById('decryptResult').innerHTML = data.result;
                    document.getElementById('decryptResultArea').style.display = 'block';
                    showStatus('decryptStatus', '✅ Text decrypted successfully!', 'success');
                }
            } catch (error) {
                showStatus('decryptStatus', '❌ Error: ' + error.message, 'error');
            }
        }
        
        function showStatus(elementId, message, type) {
            const element = document.getElementById(elementId);
            element.textContent = message;
            element.className = 'status ' + type;
            setTimeout(() => {
                element.className = 'status';
            }, 3000);
        }
        
        function copyEncryptResult() {
            if (currentEncryptResult) {
                navigator.clipboard.writeText(currentEncryptResult).then(() => {
                    showStatus('encryptStatus', '📋 Copied to clipboard!', 'success');
                }).catch(() => {
                    showStatus('encryptStatus', '❌ Failed to copy', 'error');
                });
            }
        }
        
        function copyDecryptResult() {
            if (currentDecryptResult) {
                navigator.clipboard.writeText(currentDecryptResult).then(() => {
                    showStatus('decryptStatus', '📋 Copied to clipboard!', 'success');
                }).catch(() => {
                    showStatus('decryptStatus', '❌ Failed to copy', 'error');
                });
            }
        }
        
        function clearEncrypt() {
            document.getElementById('encryptText').value = '';
            document.getElementById('encryptKey').value = '';
            document.getElementById('encryptResultArea').style.display = 'none';
            document.getElementById('encryptResult').innerHTML = '';
            currentEncryptResult = '';
        }
        
        function clearDecrypt() {
            document.getElementById('decryptText').value = '';
            document.getElementById('decryptKey').value = '';
            document.getElementById('decryptResultArea').style.display = 'none';
            document.getElementById('decryptResult').innerHTML = '';
            currentDecryptResult = '';
        }
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Error: "Invalid request"})
		return
	}
	
	encrypted, err := encrypt(req.Text, req.Key)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}
	
	json.NewEncoder(w).Encode(Response{Result: encrypted})
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Error: "Invalid request"})
		return
	}
	
	decrypted, err := decrypt(req.Text, req.Key)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}
	
	json.NewEncoder(w).Encode(Response{Result: decrypted})
}

func encrypt(plaintext, key string) (string, error) {
	if len(key) < 6 {
		return "", fmt.Errorf("key must be at least 6 characters")
	}
	
	keyHash := sha256.Sum256([]byte(key))
	
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("invalid encrypted text format")
	}
	
	keyHash := sha256.Sum256([]byte(key))
	
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("invalid encrypted text")
	}
	
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed - wrong key or corrupted data")
	}
	
	return string(plaintext), nil
}
