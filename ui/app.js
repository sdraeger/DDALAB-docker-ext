const { useState, useEffect, useRef } = React;

const API_BASE = '/api';

// Service component
function ServiceItem({ service, onAction }) {
    const [loading, setLoading] = useState(false);

    const handleAction = async (action) => {
        setLoading(true);
        await onAction(service.name, action);
        setLoading(false);
    };

    return (
        <div className="service-item">
            <div className="service-info">
                <div className="service-name">{service.name}</div>
                <div className="service-status">
                    Status: <span className={`status status-${service.status}`}>{service.status}</span>
                </div>
            </div>
            <div className="service-actions">
                {service.status === 'running' ? (
                    <>
                        <button 
                            className="btn btn-secondary" 
                            onClick={() => handleAction('restart')}
                            disabled={loading}
                        >
                            Restart
                        </button>
                        <button 
                            className="btn btn-danger" 
                            onClick={() => handleAction('stop')}
                            disabled={loading}
                        >
                            Stop
                        </button>
                    </>
                ) : (
                    <button 
                        className="btn btn-primary" 
                        onClick={() => handleAction('start')}
                        disabled={loading}
                    >
                        Start
                    </button>
                )}
                {loading && <span className="loading"></span>}
            </div>
        </div>
    );
}

// Main App component
function App() {
    const [status, setStatus] = useState({
        running: false,
        services: [],
        version: 'Unknown',
        path: 'Not found'
    });
    const [logs, setLogs] = useState('');
    const [alert, setAlert] = useState(null);
    const [loading, setLoading] = useState(false);
    const logsEndRef = useRef(null);

    useEffect(() => {
        fetchStatus();
        const interval = setInterval(fetchStatus, 5000);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        if (logsEndRef.current) {
            logsEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [logs]);

    const fetchStatus = async () => {
        try {
            const response = await fetch(`${API_BASE}/status`);
            const data = await response.json();
            setStatus(data);
        } catch (error) {
            console.error('Failed to fetch status:', error);
        }
    };

    const fetchLogs = async () => {
        try {
            const response = await fetch(`${API_BASE}/logs`);
            const data = await response.json();
            setLogs(data.logs || '');
        } catch (error) {
            console.error('Failed to fetch logs:', error);
            showAlert('Failed to fetch logs', 'error');
        }
    };

    const handleServiceAction = async (serviceName, action) => {
        try {
            const response = await fetch(`${API_BASE}/services/${serviceName}/${action}`, {
                method: 'POST'
            });
            
            if (response.ok) {
                showAlert(`${action} ${serviceName} successfully`, 'success');
                await fetchStatus();
            } else {
                throw new Error(`Failed to ${action} ${serviceName}`);
            }
        } catch (error) {
            showAlert(error.message, 'error');
        }
    };

    const handleStackAction = async (action) => {
        setLoading(true);
        try {
            const response = await fetch(`${API_BASE}/stack/${action}`, {
                method: 'POST'
            });
            
            if (response.ok) {
                showAlert(`DDALAB ${action} initiated`, 'success');
                await fetchStatus();
            } else {
                throw new Error(`Failed to ${action} DDALAB`);
            }
        } catch (error) {
            showAlert(error.message, 'error');
        } finally {
            setLoading(false);
        }
    };

    const handleBackup = async () => {
        setLoading(true);
        try {
            const response = await fetch(`${API_BASE}/backup`, {
                method: 'POST'
            });
            
            if (response.ok) {
                const data = await response.json();
                showAlert(`Backup created: ${data.filename}`, 'success');
            } else {
                throw new Error('Failed to create backup');
            }
        } catch (error) {
            showAlert(error.message, 'error');
        } finally {
            setLoading(false);
        }
    };

    const showAlert = (message, type) => {
        setAlert({ message, type });
        setTimeout(() => setAlert(null), 5000);
    };

    const runningServices = status.services.filter(s => s.status === 'running').length;

    return (
        <div className="container">
            <div className="header">
                <h1>DDALAB Manager</h1>
                <p>Manage your DDALAB installation</p>
            </div>

            {alert && (
                <div className={`alert alert-${alert.type}`}>
                    {alert.message}
                </div>
            )}

            <div className="status-grid">
                <div className="status-card">
                    <h3>Status</h3>
                    <div className="value">
                        <span className={`status status-${status.running ? 'running' : 'stopped'}`}>
                            {status.running ? 'Running' : 'Stopped'}
                        </span>
                    </div>
                </div>
                <div className="status-card">
                    <h3>Services</h3>
                    <div className="value">{runningServices} / {status.services.length}</div>
                </div>
                <div className="status-card">
                    <h3>Version</h3>
                    <div className="value">{status.version}</div>
                </div>
                <div className="status-card">
                    <h3>Installation Path</h3>
                    <div className="value" style={{ fontSize: '14px' }}>{status.path}</div>
                </div>
            </div>

            <div className="services-section">
                <h2>Services</h2>
                {status.services.map(service => (
                    <ServiceItem 
                        key={service.name} 
                        service={service} 
                        onAction={handleServiceAction}
                    />
                ))}
            </div>

            <div className="actions-section">
                <h2>Stack Actions</h2>
                <div className="action-buttons">
                    <button 
                        className="btn btn-primary" 
                        onClick={() => handleStackAction('start')}
                        disabled={loading || status.running}
                    >
                        Start All Services
                    </button>
                    <button 
                        className="btn btn-danger" 
                        onClick={() => handleStackAction('stop')}
                        disabled={loading || !status.running}
                    >
                        Stop All Services
                    </button>
                    <button 
                        className="btn btn-secondary" 
                        onClick={() => handleStackAction('restart')}
                        disabled={loading}
                    >
                        Restart All Services
                    </button>
                    <button 
                        className="btn btn-secondary" 
                        onClick={handleBackup}
                        disabled={loading}
                    >
                        Create Backup
                    </button>
                    <button 
                        className="btn btn-secondary" 
                        onClick={fetchLogs}
                        disabled={loading}
                    >
                        Refresh Logs
                    </button>
                    {loading && <span className="loading"></span>}
                </div>
            </div>

            <div className="logs-section">
                <h2>Recent Logs</h2>
                <div className="log-viewer">
                    {logs || 'No logs available. Click "Refresh Logs" to load.'}
                    <div ref={logsEndRef} />
                </div>
            </div>
        </div>
    );
}

// Render the app
ReactDOM.render(<App />, document.getElementById('root'));