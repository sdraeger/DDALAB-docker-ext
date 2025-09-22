import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  Grid,
  Card,
  CardContent,
  Chip,
  Button,
  Alert,
  CircularProgress,
  Divider,
  Stack,
  IconButton,
  Tooltip,
  AppBar,
  Toolbar,
} from '@mui/material';
import {
  PlayArrow,
  Stop,
  Refresh,
  Backup,
  Settings,
  CheckCircle,
  Error,
  Warning,
  FolderOpen,
  Dashboard,
  OpenInNew,
  Update,
  Edit,
} from '@mui/icons-material';
import PathSelector from './components/PathSelector';
import ServiceList from './components/ServiceList';
import StatusCard from './components/StatusCard';
import EnvConfigEditor from './components/EnvConfigEditor';
import { Status, Alert as AlertType, EnvConfig } from './types';

const API_BASE = 'http://ddalab-control:8080/api';

function App() {
  const [status, setStatus] = useState<Status>({
    running: false,
    services: [],
    version: 'Unknown',
    path: 'Not found'
  });
  const [currentPath, setCurrentPath] = useState<string>('');
  const [alert, setAlert] = useState<AlertType | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [initialLoading, setInitialLoading] = useState<boolean>(true);
  const [envConfig, setEnvConfig] = useState<EnvConfig | null>(null);
  const [envConfigEditorOpen, setEnvConfigEditorOpen] = useState<boolean>(false);
  const [updating, setUpdating] = useState<boolean>(false);

  useEffect(() => {
    loadInitialData();
    // Reduced polling frequency from 10s to 30s to reduce load
    const interval = setInterval(fetchStatus, 30000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (currentPath) {
      fetchStatus();
    }
  }, [currentPath]);

  const loadInitialData = async () => {
    try {
      await Promise.all([loadInitialPath(), fetchStatus(), fetchEnvConfig()]);
    } catch (error) {
      console.error('Failed to load initial data:', error);
    } finally {
      setInitialLoading(false);
    }
  };

  const loadInitialPath = async () => {
    try {
      const response = await fetch(`${API_BASE}/paths`);
      if (response.ok) {
        const data = await response.json();
        setCurrentPath(data.selected_path || '');
      }
    } catch (error) {
      console.error('Failed to load initial path:', error);
    }
  };

  const fetchStatus = async () => {
    try {
      const response = await fetch(`${API_BASE}/status`);
      const data: Status = await response.json();
      setStatus(data);
      setCurrentPath(data.path || '');
    } catch (error) {
      console.error('Failed to fetch status:', error);
      if (!alert) {
        showAlert('Failed to connect to DDALAB backend', 'error');
      }
    }
  };

  const fetchEnvConfig = async () => {
    try {
      const response = await fetch(`${API_BASE}/env`);
      if (response.ok) {
        const data: EnvConfig = await response.json();
        setEnvConfig(data);
      }
    } catch (error) {
      console.error('Failed to fetch env config:', error);
    }
  };

  const handleServiceAction = async (serviceName: string, action: string) => {
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
      showAlert((error as Error).message, 'error');
    }
  };

  const handleStackAction = async (action: string) => {
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
      showAlert((error as Error).message, 'error');
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
      showAlert((error as Error).message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (message: string, severity: AlertType['severity']) => {
    setAlert({ message, severity });
    setTimeout(() => setAlert(null), 5000);
  };

  const handlePathChange = (newPath: string) => {
    setCurrentPath(newPath);
    showAlert(`Switched to DDALAB installation at: ${newPath}`, 'success');
    // Refetch env config when path changes
    fetchEnvConfig();
  };

  const handleOpenWebsite = () => {
    if (envConfig?.url) {
      // Try Docker Desktop extension API first
      if (window.ddClient?.host?.openExternal) {
        window.ddClient.host.openExternal(envConfig.url);
        showAlert(`Opening DDALAB at ${envConfig.url}`, 'info');
      } else {
        // Fallback to window.open for development/testing
        try {
          window.open(envConfig.url, '_blank');
          showAlert(`Opening DDALAB at ${envConfig.url}`, 'info');
        } catch (error) {
          // If window.open fails, copy URL to clipboard as fallback
          navigator.clipboard.writeText(envConfig.url).then(() => {
            showAlert(`URL copied to clipboard: ${envConfig.url}`, 'info');
          }).catch(() => {
            showAlert(`Please open: ${envConfig.url}`, 'info');
          });
        }
      }
    } else {
      showAlert('DDALAB URL not configured or found', 'warning');
    }
  };

  const handleUpdateDDALAB = async () => {
    setUpdating(true);
    try {
      const response = await fetch(`${API_BASE}/update`, {
        method: 'POST'
      });
      
      if (response.ok) {
        const data = await response.json();
        showAlert(data.message || 'DDALAB updated successfully', 'success');
        await fetchStatus();
      } else {
        throw new Error('Failed to update DDALAB');
      }
    } catch (error) {
      showAlert((error as Error).message, 'error');
    } finally {
      setUpdating(false);
    }
  };

  const handleEnvConfigSave = () => {
    showAlert('Environment configuration saved successfully', 'success');
    fetchEnvConfig(); // Refresh env config
  };

  const runningServices = status.services.filter(s => s.status === 'running').length;
  const overallStatus = status.running ? 'running' : 'stopped';

  if (initialLoading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4, display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
        <Stack spacing={2} alignItems="center">
          <CircularProgress size={48} />
          <Typography variant="h6" color="text.secondary">
            Loading DDALAB Manager...
          </Typography>
        </Stack>
      </Container>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* Docker-styled Header */}
      <AppBar position="static" elevation={0} sx={{ bgcolor: 'primary.main' }}>
        <Toolbar>
          <Dashboard sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            DDALAB Manager
          </Typography>
          <Tooltip title="Refresh Status">
            <IconButton color="inherit" onClick={fetchStatus}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 3 }}>
        {/* Alert */}
        {alert && (
          <Alert 
            severity={alert.severity} 
            sx={{ mb: 3 }} 
            onClose={() => setAlert(null)}
          >
            {alert.message}
          </Alert>
        )}

        {/* Path Selector */}
        <Paper elevation={1} sx={{ p: 3, mb: 3 }}>
          <PathSelector
            currentPath={currentPath}
            onPathChange={handlePathChange}
          />
        </Paper>

        {/* Status Overview */}
        <Grid container spacing={3} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard
              title="Overall Status"
              value={overallStatus}
              type="status"
              icon={overallStatus === 'running' ? <CheckCircle color="success" /> : <Error color="error" />}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard
              title="Services"
              value={`${runningServices} / ${status.services.length}`}
              subtitle="Running"
              icon={<Settings color="primary" />}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard
              title="Version"
              value={status.version}
              icon={<FolderOpen color="info" />}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard
              title="Health"
              value={status.running ? "Healthy" : "Stopped"}
              type="health"
              icon={status.running ? <CheckCircle color="success" /> : <Warning color="warning" />}
            />
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          {/* Stack Actions */}
          <Grid item xs={12} md={6}>
            <Card elevation={1}>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <PlayArrow color="primary" />
                  Stack Actions
                </Typography>
                <Divider sx={{ mb: 2 }} />
                
                <Stack spacing={2}>
                  {status.running ? (
                    <>
                      <Button
                        variant="outlined"
                        startIcon={<Refresh />}
                        onClick={() => handleStackAction('restart')}
                        disabled={loading || !currentPath}
                        fullWidth
                      >
                        Restart All Services
                      </Button>
                      <Button
                        variant="outlined"
                        color="error"
                        startIcon={<Stop />}
                        onClick={() => handleStackAction('stop')}
                        disabled={loading || !currentPath}
                        fullWidth
                      >
                        Stop All Services
                      </Button>
                    </>
                  ) : (
                    <Button
                      variant="contained"
                      startIcon={<PlayArrow />}
                      onClick={() => handleStackAction('start')}
                      disabled={loading || !currentPath}
                      fullWidth
                      size="large"
                    >
                      Start All Services
                    </Button>
                  )}
                  
                  <Button
                    variant="outlined"
                    startIcon={<Backup />}
                    onClick={handleBackup}
                    disabled={loading || !status.running || !currentPath}
                    fullWidth
                  >
                    Backup Database
                  </Button>
                  
                  <Button
                    variant="outlined"
                    startIcon={<Edit />}
                    onClick={() => setEnvConfigEditorOpen(true)}
                    disabled={!currentPath}
                    fullWidth
                  >
                    Edit Configuration
                  </Button>
                  
                  <Button
                    variant="outlined"
                    startIcon={<Update />}
                    onClick={handleUpdateDDALAB}
                    disabled={updating || !currentPath}
                    fullWidth
                  >
                    {updating ? 'Updating...' : 'Update DDALAB'}
                  </Button>
                  
                  <Button
                    variant="outlined"
                    color="primary"
                    startIcon={<OpenInNew />}
                    onClick={handleOpenWebsite}
                    disabled={!envConfig?.url || !currentPath}
                    fullWidth
                  >
                    Open DDALAB Website
                  </Button>

                  {!currentPath && (
                    <Alert severity="warning" sx={{ mt: 1 }}>
                      Select a DDALAB installation path to control services
                    </Alert>
                  )}
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          {/* Services */}
          <Grid item xs={12} md={6}>
            <Card elevation={1}>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Settings color="primary" />
                  Services
                </Typography>
                <Divider sx={{ mb: 2 }} />
                
                <ServiceList
                  services={status.services}
                  onServiceAction={handleServiceAction}
                  disabled={!currentPath}
                />
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {loading && (
          <Box 
            position="fixed" 
            top={0} 
            left={0} 
            right={0} 
            bottom={0} 
            display="flex" 
            alignItems="center" 
            justifyContent="center" 
            bgcolor="rgba(0,0,0,0.5)"
            zIndex={9999}
          >
            <Paper elevation={8} sx={{ p: 3 }}>
              <Stack direction="row" spacing={2} alignItems="center">
                <CircularProgress />
                <Typography>Processing...</Typography>
              </Stack>
            </Paper>
          </Box>
        )}
      </Container>

      <EnvConfigEditor
        open={envConfigEditorOpen}
        onClose={() => setEnvConfigEditorOpen(false)}
        onSave={handleEnvConfigSave}
      />
    </Box>
  );
}

export default App;