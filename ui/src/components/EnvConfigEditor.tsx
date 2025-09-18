import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  InputAdornment,
  IconButton,
  Box,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Grid,
  Chip,
  Alert,
  CircularProgress,
  Tooltip,
  FormControlLabel,
  Switch,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Checkbox,
} from '@mui/material';
import {
  Search,
  ExpandMore,
  Visibility,
  VisibilityOff,
  Save,
  Close,
  Warning,
  Check,
  FilterList,
  Sort,
  Download,
  Upload,
  Backup,
  Restore,
} from '@mui/icons-material';
import { EnvFile, EnvVar, ValidationResult, ValidationError } from '../types';

interface EnvConfigEditorProps {
  open: boolean;
  onClose: () => void;
  onSave: () => void;
}

const API_BASE = 'http://localhost:8080/api';

const EnvConfigEditor: React.FC<EnvConfigEditorProps> = ({ open, onClose, onSave }) => {
  const [envFile, setEnvFile] = useState<EnvFile | null>(null);
  const [editedVars, setEditedVars] = useState<Map<string, string>>(new Map());
  const [searchTerm, setSearchTerm] = useState('');
  const [showSecrets, setShowSecrets] = useState(false);
  const [filterRequired, setFilterRequired] = useState(false);
  const [filterSecret, setFilterSecret] = useState(false);
  const [sortBy, setSortBy] = useState<'name' | 'section' | 'required'>('section');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [validation, setValidation] = useState<ValidationResult | null>(null);
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set());

  useEffect(() => {
    if (open) {
      loadEnvFile();
    }
  }, [open]);

  const loadEnvFile = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/env/file`);
      if (response.ok) {
        const data: EnvFile = await response.json();
        setEnvFile(data);
        // Expand first section by default
        const sections = [...new Set(data.variables.map(v => v.section))];
        if (sections.length > 0) {
          setExpandedSections(new Set([sections[0]]));
        }
      } else {
        console.error('Failed to load env file');
      }
    } catch (error) {
      console.error('Error loading env file:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!envFile) return;

    setSaving(true);
    try {
      // Update variables with edited values
      const updatedVars = envFile.variables.map(v => ({
        ...v,
        value: editedVars.has(v.key) ? editedVars.get(v.key)! : v.value,
      }));

      const response = await fetch(`${API_BASE}/env/file`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ variables: updatedVars }),
      });

      if (response.ok) {
        const result: ValidationResult = await response.json();
        if (result.valid) {
          onSave();
          onClose();
        } else {
          setValidation(result);
        }
      } else {
        console.error('Failed to save env file');
      }
    } catch (error) {
      console.error('Error saving env file:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleValidate = async () => {
    if (!envFile) return;

    try {
      const updatedVars = envFile.variables.map(v => ({
        ...v,
        value: editedVars.has(v.key) ? editedVars.get(v.key)! : v.value,
      }));

      const response = await fetch(`${API_BASE}/env/validate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ variables: updatedVars }),
      });

      if (response.ok) {
        const result: ValidationResult = await response.json();
        setValidation(result);
      }
    } catch (error) {
      console.error('Error validating env file:', error);
    }
  };

  const handleVarChange = (key: string, value: string) => {
    const newEditedVars = new Map(editedVars);
    newEditedVars.set(key, value);
    setEditedVars(newEditedVars);
    setValidation(null); // Clear validation when editing
  };

  const toggleSection = (section: string) => {
    const newExpanded = new Set(expandedSections);
    if (newExpanded.has(section)) {
      newExpanded.delete(section);
    } else {
      newExpanded.add(section);
    }
    setExpandedSections(newExpanded);
  };

  const getFilteredVars = () => {
    if (!envFile) return [];
    
    let filtered = envFile.variables.filter(v => {
      // Search filter
      if (searchTerm) {
        const searchLower = searchTerm.toLowerCase();
        const matchesSearch = (
          v.key.toLowerCase().includes(searchLower) ||
          v.value.toLowerCase().includes(searchLower) ||
          v.comment.toLowerCase().includes(searchLower) ||
          v.section.toLowerCase().includes(searchLower)
        );
        if (!matchesSearch) return false;
      }

      // Required filter
      if (filterRequired && !v.required) return false;

      // Secret filter
      if (filterSecret && !v.secret) return false;

      return true;
    });

    // Sort the results
    filtered.sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return a.key.localeCompare(b.key);
        case 'required':
          if (a.required !== b.required) {
            return a.required ? -1 : 1;
          }
          return a.key.localeCompare(b.key);
        case 'section':
        default:
          if (a.section !== b.section) {
            return a.section.localeCompare(b.section);
          }
          return a.key.localeCompare(b.key);
      }
    });

    return filtered;
  };

  const getVarsBySection = () => {
    const vars = getFilteredVars();
    const sections = new Map<string, EnvVar[]>();
    
    vars.forEach(v => {
      if (!sections.has(v.section)) {
        sections.set(v.section, []);
      }
      sections.get(v.section)!.push(v);
    });
    
    return sections;
  };

  const getValidationError = (key: string): ValidationError | undefined => {
    if (!validation || validation.valid) return undefined;
    return validation.errors.find(e => e.key === key);
  };

  const handleExportConfig = () => {
    if (!envFile) return;

    const exportData = {
      exported_at: new Date().toISOString(),
      variables: envFile.variables.map(v => ({
        key: v.key,
        value: v.secret ? '***REDACTED***' : v.value,
        comment: v.comment,
        section: v.section,
        required: v.required,
        secret: v.secret,
      })),
    };

    const blob = new Blob([JSON.stringify(exportData, null, 2)], {
      type: 'application/json',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `ddalab-config-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleImportConfig = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const content = e.target?.result as string;
        const importData = JSON.parse(content);
        
        if (importData.variables && Array.isArray(importData.variables)) {
          // Update edited vars with imported values (excluding secrets)
          const newEditedVars = new Map(editedVars);
          importData.variables.forEach((imported: any) => {
            if (imported.key && !imported.secret && imported.value !== '***REDACTED***') {
              newEditedVars.set(imported.key, imported.value);
            }
          });
          setEditedVars(newEditedVars);
        }
      } catch (error) {
        console.error('Failed to import configuration:', error);
        // You could add a proper error notification here
      }
    };
    reader.readAsText(file);
  };

  const renderEnvVar = (envVar: EnvVar) => {
    const value = editedVars.has(envVar.key) ? editedVars.get(envVar.key)! : envVar.value;
    const error = getValidationError(envVar.key);
    const isSecret = envVar.secret && !showSecrets;

    return (
      <Box key={envVar.key} mb={2}>
        <TextField
          fullWidth
          label={
            <Box display="flex" alignItems="center" gap={1}>
              {envVar.key}
              {envVar.required && <Chip label="Required" size="small" color="error" />}
              {envVar.secret && <Chip label="Secret" size="small" color="warning" />}
            </Box>
          }
          value={isSecret ? '••••••••' : value}
          onChange={(e) => handleVarChange(envVar.key, e.target.value)}
          error={!!error}
          helperText={error ? error.message : envVar.comment}
          disabled={isSecret}
          InputProps={{
            endAdornment: envVar.secret && (
              <InputAdornment position="end">
                <IconButton
                  onClick={() => setShowSecrets(!showSecrets)}
                  edge="end"
                  size="small"
                >
                  {showSecrets ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
      </Box>
    );
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">Environment Configuration</Typography>
          <IconButton onClick={onClose} size="small">
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent>
        {loading ? (
          <Box display="flex" justifyContent="center" p={4}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Box mb={2}>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={4}>
                  <TextField
                    fullWidth
                    placeholder="Search variables..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <Search />
                        </InputAdornment>
                      ),
                    }}
                  />
                </Grid>
                <Grid item xs={12} sm={3}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Sort by</InputLabel>
                    <Select
                      value={sortBy}
                      label="Sort by"
                      onChange={(e) => setSortBy(e.target.value as any)}
                      startAdornment={<Sort sx={{ mr: 1, color: 'action.active' }} />}
                    >
                      <MenuItem value="section">Section</MenuItem>
                      <MenuItem value="name">Name</MenuItem>
                      <MenuItem value="required">Required First</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={12} sm={5}>
                  <Box display="flex" gap={1} alignItems="center">
                    <FormControlLabel
                      control={
                        <Switch
                          checked={showSecrets}
                          onChange={(e) => setShowSecrets(e.target.checked)}
                          size="small"
                        />
                      }
                      label="Show secrets"
                    />
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={filterRequired}
                          onChange={(e) => setFilterRequired(e.target.checked)}
                          size="small"
                        />
                      }
                      label="Required only"
                    />
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={filterSecret}
                          onChange={(e) => setFilterSecret(e.target.checked)}
                          size="small"
                        />
                      }
                      label="Secrets only"
                    />
                  </Box>
                </Grid>
              </Grid>
            </Box>

            {validation && !validation.valid && (
              <Alert severity="error" sx={{ mb: 2 }}>
                <Typography variant="subtitle2" gutterBottom>
                  Validation errors found:
                </Typography>
                {validation.errors.map((err) => (
                  <Typography key={err.key} variant="body2">
                    • {err.key}: {err.message}
                  </Typography>
                ))}
              </Alert>
            )}

            {validation && validation.valid && (
              <Alert severity="success" sx={{ mb: 2 }}>
                Configuration is valid
              </Alert>
            )}

            {Array.from(getVarsBySection().entries()).map(([section, vars]) => (
              <Accordion
                key={section}
                expanded={expandedSections.has(section)}
                onChange={() => toggleSection(section)}
              >
                <AccordionSummary expandIcon={<ExpandMore />}>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Typography variant="subtitle1" fontWeight="bold">
                      {section}
                    </Typography>
                    <Chip label={vars.length} size="small" />
                  </Box>
                </AccordionSummary>
                <AccordionDetails>
                  {vars.map(renderEnvVar)}
                </AccordionDetails>
              </Accordion>
            ))}
          </>
        )}
      </DialogContent>

      <DialogActions>
        <Box display="flex" gap={1}>
          <Button onClick={handleValidate} startIcon={<Check />}>
            Validate
          </Button>
          <Button onClick={handleExportConfig} startIcon={<Download />}>
            Export
          </Button>
          <Button component="label" startIcon={<Upload />}>
            Import
            <input
              type="file"
              accept=".json"
              hidden
              onChange={handleImportConfig}
            />
          </Button>
        </Box>
        <Box flexGrow={1} />
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={handleSave}
          variant="contained"
          startIcon={<Save />}
          disabled={saving || loading}
        >
          {saving ? 'Saving...' : 'Save Changes'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default EnvConfigEditor;