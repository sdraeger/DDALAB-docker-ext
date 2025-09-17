import React, { useState } from 'react';
import {
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Button,
  ButtonGroup,
  Chip,
  CircularProgress,
  Box,
  Typography,
  Stack,
} from '@mui/material';
import {
  PlayArrow,
  Stop,
  Refresh,
} from '@mui/icons-material';
import { Service, ServiceListProps } from '../types';

interface ServiceItemProps {
  service: Service;
  onAction: (serviceName: string, action: string) => Promise<void>;
  disabled?: boolean;
}

function ServiceItem({ service, onAction, disabled }: ServiceItemProps) {
  const [loading, setLoading] = useState<boolean>(false);

  const handleAction = async (action: string) => {
    setLoading(true);
    try {
      await onAction(service.name, action);
    } catch (error) {
      console.error('Action failed:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string): 'success' | 'error' | 'default' => {
    switch (status) {
      case 'running':
        return 'success';
      case 'stopped':
        return 'error';
      default:
        return 'default';
    }
  };

  return (
    <ListItem divider>
      <ListItemText
        primary={
          <Stack direction="row" spacing={1} alignItems="center">
            <Typography variant="body1" fontWeight="medium">
              {service.name}
            </Typography>
            <Chip
              label={service.status}
              color={getStatusColor(service.status)}
              size="small"
              variant="outlined"
            />
          </Stack>
        }
        secondary={`Status: ${service.status}`}
      />
      <ListItemSecondaryAction>
        <Stack direction="row" spacing={1} alignItems="center">
          {service.status === 'running' ? (
            <ButtonGroup size="small" variant="outlined">
              <Button
                startIcon={<Refresh />}
                onClick={() => handleAction('restart')}
                disabled={loading || disabled}
              >
                Restart
              </Button>
              <Button
                startIcon={<Stop />}
                color="error"
                onClick={() => handleAction('stop')}
                disabled={loading || disabled}
              >
                Stop
              </Button>
            </ButtonGroup>
          ) : (
            <Button
              variant="contained"
              startIcon={<PlayArrow />}
              onClick={() => handleAction('start')}
              disabled={loading || disabled}
              size="small"
            >
              Start
            </Button>
          )}
          {loading && <CircularProgress size={20} sx={{ ml: 1 }} />}
        </Stack>
      </ListItemSecondaryAction>
    </ListItem>
  );
}

function ServiceList({ services, onServiceAction, disabled }: ServiceListProps) {
  if (!services || services.length === 0) {
    return (
      <Box textAlign="center" py={3}>
        <Typography variant="body2" color="text.secondary">
          No services found
        </Typography>
      </Box>
    );
  }

  return (
    <List disablePadding>
      {services.map((service) => (
        <ServiceItem
          key={service.name}
          service={service}
          onAction={onServiceAction}
          disabled={disabled}
        />
      ))}
    </List>
  );
}

export default ServiceList;