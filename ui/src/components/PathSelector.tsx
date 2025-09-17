import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Alert,
  CircularProgress,
  Chip,
  Paper,
  IconButton,
  Tooltip,
  Divider,
  Stack,
} from '@mui/material';
import {
  FolderOpen,
  Check,
  Close,
  Refresh,
  Settings,
  CheckCircle,
  Error,
} from '@mui/icons-material';
import { PathSelectorProps, PathValidationResult, ExtensionConfig } from '../types';

const API_BASE = 'http://localhost:8080/api';

function PathSelector({ currentPath, onPathChange }: PathSelectorProps) {
  const [open, setOpen] = useState<boolean>(false);
  const [customPath, setCustomPath] = useState<string>('');
  const [discoveredPaths, setDiscoveredPaths] = useState<string[]>([]);
  const [knownPaths, setKnownPaths] = useState<string[]>([]);
  const [validationResult, setValidationResult] = useState<PathValidationResult | null>(null);
  const [isValidating, setIsValidating] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (open) {
      loadPaths();
    }
  }, [open]);

  const loadPaths = async () => {
    setLoading(true);
    try {
      // Load known paths
      const pathsResponse = await fetch(`${API_BASE}/paths`);
      if (pathsResponse.ok) {
        const pathsData: ExtensionConfig = await pathsResponse.json();
        setKnownPaths(pathsData.known_paths || []);
      }

      // Discover new paths
      const discoverResponse = await fetch(`${API_BASE}/paths/discover`);
      if (discoverResponse.ok) {
        const discoverData: { discovered_paths: string[] } = await discoverResponse.json();
        setDiscoveredPaths(discoverData.discovered_paths || []);
      }
    } catch (error) {
      console.error('Failed to load paths:', error);
    } finally {
      setLoading(false);
    }
  };

  const validatePath = async (path: string) => {
    if (!path.trim()) {
      setValidationResult(null);
      return;
    }

    setIsValidating(true);
    try {
      const response = await fetch(`${API_BASE}/paths/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path: path.trim() }),
      });

      const result: PathValidationResult = await response.json();
      setValidationResult(result);
    } catch (error) {
      setValidationResult({
        valid: false,
        path: path,
        message: 'Failed to validate path: ' + (error as Error).message,
        has_compose: false,
        has_ddalab_script: false,
      });
    } finally {
      setIsValidating(false);
    }
  };

  const selectPath = async (path: string) => {
    try {
      const response = await fetch(`${API_BASE}/paths/select`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path }),
      });

      const result: PathValidationResult = await response.json();
      if (response.ok && result.valid) {
        onPathChange(path);
        setOpen(false);
        setCustomPath('');
        setValidationResult(null);
      } else {
        setValidationResult(result);
      }
    } catch (error) {
      setValidationResult({
        valid: false,
        path: path,
        message: 'Failed to select path: ' + (error as Error).message,
        has_compose: false,
        has_ddalab_script: false,
      });
    }
  };

  const handleCustomPathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const path = e.target.value;
    setCustomPath(path);
    validatePath(path);
  };

  const handleClose = () => {
    setOpen(false);
    setCustomPath('');
    setValidationResult(null);
  };

  const allPaths = [...new Set([...knownPaths, ...discoveredPaths])];

  return (
    <>
      {/* Current Path Display */}
      <Box>
        <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 1 }}>
          <FolderOpen color="primary" />
          <Typography variant="h6" component="div">
            Installation Path
          </Typography>
        </Stack>
        <Stack direction="row" spacing={2} alignItems="center">
          <Paper 
            elevation={1}
            sx={{ 
              flex: 1,
              fontFamily: 'monospace !important',
              wordBreak: 'break-all',
              bgcolor: '#2e2e2e !important',
              color: '#ffffff !important',
              p: 1.5,
              borderRadius: 1,
              border: '2px solid',
              borderColor: '#555555 !important',
              fontSize: '0.875rem',
              minHeight: '20px',
              display: 'flex',
              alignItems: 'center'
            }}
          >
            <Typography 
              variant="body2"
              sx={{ 
                fontFamily: 'monospace !important',
                color: '#ffffff !important',
                fontSize: 'inherit'
              }}
            >
              {currentPath || 'Not selected'}
            </Typography>
          </Paper>
          <Button
            variant="outlined"
            startIcon={<Settings />}
            onClick={() => setOpen(true)}
            size="large"
          >
            Change Path
          </Button>
        </Stack>
      </Box>

      {/* Path Selection Dialog */}
      <Dialog 
        open={open} 
        onClose={handleClose} 
        maxWidth="md" 
        fullWidth
        PaperProps={{
          sx: { minHeight: '500px' }
        }}
      >
        <DialogTitle>
          <Stack direction="row" spacing={1} alignItems="center">
            <FolderOpen color="primary" />
            <Typography variant="h6">
              Select DDALAB Installation Path
            </Typography>
          </Stack>
        </DialogTitle>
        
        <DialogContent dividers>
          {/* Custom Path Input */}
          <Box mb={3}>
            <Typography variant="h6" gutterBottom>
              Custom Path
            </Typography>
            <TextField
              fullWidth
              label="Enter path to DDALAB installation"
              placeholder="/Users/username/DDALAB-setup"
              value={customPath}
              onChange={handleCustomPathChange}
              variant="outlined"
              InputProps={{
                style: { 
                  fontFamily: 'monospace',
                },
                endAdornment: isValidating ? <CircularProgress size={20} /> : undefined
              }}
              sx={{
                // Force text visibility in Docker's dark mode
                '& input, & input[type="text"], & textarea': {
                  color: '#ffffff !important',
                  fontFamily: 'monospace !important',
                  WebkitTextFillColor: '#ffffff !important',
                  caretColor: '#ffffff !important',
                  '&::placeholder': {
                    color: '#aaaaaa !important',
                    opacity: '1 !important',
                  },
                  '&:focus': {
                    color: '#ffffff !important',
                    WebkitTextFillColor: '#ffffff !important',
                  },
                  '&:-webkit-autofill': {
                    WebkitBoxShadow: '0 0 0 1000px #2e2e2e inset !important',
                    WebkitTextFillColor: '#ffffff !important',
                    transition: 'background-color 5000s ease-in-out 0s',
                  },
                  '&:-webkit-autofill:hover': {
                    WebkitBoxShadow: '0 0 0 1000px #2e2e2e inset !important',
                    WebkitTextFillColor: '#ffffff !important',
                  },
                  '&:-webkit-autofill:focus': {
                    WebkitBoxShadow: '0 0 0 1000px #2e2e2e inset !important',
                    WebkitTextFillColor: '#ffffff !important',
                  },
                },
                '& .MuiInputBase-input': {
                  color: '#ffffff !important',
                  fontFamily: 'monospace !important',
                  WebkitTextFillColor: '#ffffff !important',
                  caretColor: '#ffffff !important',
                  '&::selection': {
                    backgroundColor: '#ffffff !important',
                    color: '#2e2e2e !important',
                  },
                  '&::-moz-selection': {
                    backgroundColor: '#ffffff !important',
                    color: '#2e2e2e !important',
                  },
                },
                '& .MuiInputBase-root': {
                  bgcolor: '#2e2e2e !important',
                  color: '#ffffff !important',
                  '&:hover': {
                    bgcolor: '#3a3a3a !important',
                  },
                  '&.Mui-focused': {
                    bgcolor: '#3a3a3a !important',
                  }
                },
                '& .MuiOutlinedInput-root': {
                  color: '#ffffff !important',
                  backgroundColor: '#2e2e2e !important',
                  '& fieldset': {
                    borderColor: '#666666 !important',
                    borderWidth: '2px !important',
                  },
                  '&:hover fieldset': {
                    borderColor: '#888888 !important',
                  },
                  '&.Mui-focused fieldset': {
                    borderColor: '#1976d2 !important',
                  },
                },
                '& .MuiInputLabel-root': {
                  color: '#cccccc !important',
                  '&.Mui-focused': {
                    color: '#1976d2 !important',
                  }
                }
              }}
            />
            
            {/* Validation Result */}
            {validationResult && (
              <Alert 
                severity={validationResult.valid ? 'success' : 'error'} 
                sx={{ mt: 2 }}
                icon={validationResult.valid ? <CheckCircle /> : <Error />}
              >
                <Typography variant="body2" fontWeight="medium">
                  {validationResult.message}
                </Typography>
                {validationResult.valid && (
                  <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                    {validationResult.has_compose && (
                      <Chip label="docker-compose.yml" size="small" color="success" variant="outlined" />
                    )}
                    {validationResult.has_ddalab_script && (
                      <Chip label="DDALAB script" size="small" color="success" variant="outlined" />
                    )}
                  </Stack>
                )}
              </Alert>
            )}

            {validationResult && validationResult.valid && (
              <Box mt={2}>
                <Button
                  variant="contained"
                  startIcon={<Check />}
                  onClick={() => selectPath(customPath)}
                  fullWidth
                >
                  Select This Path
                </Button>
              </Box>
            )}
          </Box>

          <Divider sx={{ my: 3 }} />

          {/* Discovered Paths */}
          {allPaths.length > 0 && (
            <Box>
              <Stack direction="row" alignItems="center" justifyContent="space-between" mb={2}>
                <Typography variant="h6">
                  Discovered Installations
                </Typography>
                <Tooltip title="Refresh">
                  <IconButton onClick={loadPaths} disabled={loading}>
                    <Refresh />
                  </IconButton>
                </Tooltip>
              </Stack>
              
              {loading ? (
                <Box textAlign="center" py={3}>
                  <CircularProgress />
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    Searching for DDALAB installations...
                  </Typography>
                </Box>
              ) : (
                <Paper variant="outlined">
                  <List disablePadding>
                    {allPaths.map((path, index) => (
                      <ListItem key={index} divider={index < allPaths.length - 1}>
                        <ListItemText
                          primary={
                            <Paper
                              elevation={0}
                              sx={{ 
                                fontFamily: 'monospace',
                                wordBreak: 'break-all',
                                color: (theme) => theme.palette.mode === 'dark' ? '#ffffff !important' : '#000000 !important',
                                bgcolor: (theme) => theme.palette.mode === 'dark' ? '#2e2e2e !important' : '#f9f9f9 !important',
                                p: 1,
                                borderRadius: 0.5,
                                fontSize: '0.875rem',
                                border: '1px solid',
                                borderColor: (theme) => theme.palette.mode === 'dark' ? '#555555' : '#e0e0e0'
                              }}
                            >
                              <Typography
                                variant="body1"
                                sx={{ 
                                  fontFamily: 'monospace',
                                  color: 'inherit',
                                  fontSize: 'inherit'
                                }}
                              >
                                {path}
                              </Typography>
                            </Paper>
                          }
                          secondary={
                            knownPaths.includes(path) ? (
                              <Chip 
                                label="Previously used" 
                                size="small" 
                                color="primary" 
                                variant="outlined"
                                sx={{ mt: 1 }}
                              />
                            ) : null
                          }
                        />
                        <ListItemSecondaryAction>
                          <Button
                            variant="contained"
                            size="small"
                            onClick={() => selectPath(path)}
                          >
                            Select
                          </Button>
                        </ListItemSecondaryAction>
                      </ListItem>
                    ))}
                  </List>
                </Paper>
              )}
            </Box>
          )}

          <Divider sx={{ my: 3 }} />

          {/* Help Section */}
          <Box>
            <Typography variant="h6" gutterBottom>
              Help
            </Typography>
            <Alert severity="info">
              <Typography variant="body2" paragraph>
                A valid DDALAB installation should contain:
              </Typography>
              <Box component="ul" sx={{ pl: 2, mb: 1 }}>
                <li><code>docker-compose.yml</code> - Main deployment configuration</li>
                <li><code>ddalab.sh</code> or <code>ddalab.ps1</code> - Control script</li>
                <li>DDALAB services configuration</li>
              </Box>
              <Typography variant="body2" paragraph>
                Common locations:
              </Typography>
              <Box component="ul" sx={{ pl: 2, mb: 0 }}>
                <li><code>~/DDALAB-setup</code></li>
                <li><code>~/Desktop/DDALAB-setup</code></li>
                <li><code>~/Documents/DDALAB-setup</code></li>
              </Box>
            </Alert>
          </Box>
        </DialogContent>
        
        <DialogActions>
          <Button onClick={handleClose} startIcon={<Close />}>
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}

export default PathSelector;