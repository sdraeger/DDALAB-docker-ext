import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
} from '@mui/material';
import { StatusCardProps } from '../types';

function StatusCard({ title, value, subtitle, type = 'default', icon }: StatusCardProps) {
  const getStatusColor = (): 'success' | 'error' | 'warning' | 'primary' => {
    if (type === 'status') {
      return value === 'running' ? 'success' : 'error';
    }
    if (type === 'health') {
      return value === 'Healthy' ? 'success' : 'warning';
    }
    return 'primary';
  };

  const getStatusVariant = (): 'filled' | 'outlined' => {
    if (type === 'status' || type === 'health') {
      return 'filled';
    }
    return 'outlined';
  };

  return (
    <Card elevation={1} sx={{ height: '100%' }}>
      <CardContent>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
          <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 500, textTransform: 'uppercase', letterSpacing: '0.5px' }}>
            {title}
          </Typography>
          {icon}
        </Box>
        
        <Box>
          {(type === 'status' || type === 'health') ? (
            <Chip
              label={value}
              color={getStatusColor()}
              variant={getStatusVariant()}
              size="small"
              sx={{ fontWeight: 'bold' }}
            />
          ) : (
            <Typography variant="h5" component="div" fontWeight="bold">
              {value}
            </Typography>
          )}
          
          {subtitle && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              {subtitle}
            </Typography>
          )}
        </Box>
      </CardContent>
    </Card>
  );
}

export default StatusCard;