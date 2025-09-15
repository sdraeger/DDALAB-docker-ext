const { useState, useEffect } = React;

// Simple App component without backend
function App() {
    const [setupPath, setSetupPath] = useState('');

    return (
        <div className="container">
            <div className="header">
                <h1>DDALAB Manager</h1>
                <p>Manage your DDALAB installation</p>
            </div>

            <div className="status-card" style={{ marginBottom: '30px' }}>
                <h3>Getting Started</h3>
                <p style={{ marginTop: '10px' }}>
                    This extension helps you manage DDALAB installations. 
                    Make sure you have DDALAB-setup installed on your system.
                </p>
            </div>

            <div className="services-section">
                <h2>Quick Actions</h2>
                <p>To manage DDALAB, use these commands in your terminal:</p>
                <div style={{ 
                    backgroundColor: '#1e293b', 
                    color: '#e2e8f0', 
                    padding: '15px', 
                    borderRadius: '6px', 
                    fontFamily: 'monospace',
                    marginTop: '15px'
                }}>
                    <div>cd ~/DDALAB-setup</div>
                    <div>./ddalab.sh start    # Start DDALAB</div>
                    <div>./ddalab.sh stop     # Stop DDALAB</div>
                    <div>./ddalab.sh status   # Check status</div>
                    <div>./ddalab.sh logs     # View logs</div>
                </div>
            </div>

            <div className="actions-section" style={{ marginTop: '30px' }}>
                <h2>DDALAB Setup Path</h2>
                <input 
                    type="text" 
                    placeholder="Enter path to DDALAB-setup directory"
                    value={setupPath}
                    onChange={(e) => setSetupPath(e.target.value)}
                    style={{
                        width: '100%',
                        padding: '10px',
                        borderRadius: '6px',
                        border: '1px solid #e2e8f0',
                        marginTop: '10px'
                    }}
                />
                <p style={{ marginTop: '10px', fontSize: '14px', color: '#64748b' }}>
                    Common locations: ~/DDALAB-setup, ~/Desktop/DDALAB-setup
                </p>
            </div>
        </div>
    );
}

// Render the app
ReactDOM.render(<App />, document.getElementById('root'));