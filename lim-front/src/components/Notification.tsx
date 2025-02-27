import React from 'react';
import { Snackbar, Alert, SnackbarOrigin, AlertColor } from '@mui/material';

interface NotificationProps {
  open: boolean;
  message: string | null;
  severity?: AlertColor; // 'error' | 'warning' | 'info' | 'success'
  autoHideDuration?: number;
  onClose: () => void;
  anchorOrigin?: SnackbarOrigin;
}

const Notification: React.FC<NotificationProps> = ({
  open,
  message,
  severity = 'info',
  autoHideDuration = 3000,
  onClose,
  anchorOrigin = { vertical: 'top', horizontal: 'center' },
}) => {
  return (
    <Snackbar
      open={open}
      autoHideDuration={autoHideDuration}
      onClose={onClose}
      anchorOrigin={anchorOrigin}
    >
      <Alert severity={severity} onClose={onClose}>
        {message}
      </Alert>
    </Snackbar>
  );
};

export default Notification;
