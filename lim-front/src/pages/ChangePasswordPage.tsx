import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ClientService } from '../api/services/clientService';
import { Button, Box, Typography } from '@mui/material';
import Notification from '../components/Notification';
import InputField from '../components/InputField';

const ChangePasswordPage: React.FC = () => {
  const navigate = useNavigate();
  const [newPassword, setNewPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [openSnackbar, setOpenSnackbar] = useState(false);

  const clientService = new ClientService();

  const handleChangePassword = async () => {
    if (newPassword !== confirmPassword) {
      setError('Пароли не совпадают');
      setOpenSnackbar(true);
      return;
    }

    try {
      await clientService.changePassword(newPassword, setError);
      setSuccessMessage('Пароль успешно изменен');
      setError(null);
      setOpenSnackbar(true);
    } catch (err: any) {
      setSuccessMessage(null);
      setError(err);
      setOpenSnackbar(true);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        backgroundColor: 'background.default',
        padding: 3,
      }}
    >
      <Button
        variant="text"
        color="inherit"
        onClick={() => navigate('/schedule')}
        sx={{
          position: 'absolute',
          top: 16,
          left: 16,
        }}
      >
        Домой
      </Button>
      <Typography variant="h1" sx={{ marginBottom: 3 }}>
        Смена пароля
      </Typography>

      <form style={{ width: '100%', maxWidth: '400px' }} noValidate autoComplete="off">
        <InputField
          label="Новый пароль"
          name="newPassword"
          type="password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
        />
        <InputField
          label="Подтвердите новый пароль"
          name="confirmPassword"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
        />
        <Box sx={{ display: 'flex', justifyContent: 'space-between', marginBottom: 2 }}>
          <Button
            type="button"
            variant="contained"
            color="primary"
            onClick={handleChangePassword}
            sx={{ width: '48%' }}
          >
            Изменить пароль
          </Button>
          <Button
            variant="text"
            color="inherit"
            onClick={() => navigate('/schedule')}
            sx={{ width: '48%' }}
          >
            Отмена
          </Button>
        </Box>
      </form>

      <Notification
        open={openSnackbar}
        message={error || successMessage}
        severity={error ? 'error' : 'success'}
        onClose={() => setOpenSnackbar(false)}
      />
    </Box>
  );
};

export default ChangePasswordPage;
