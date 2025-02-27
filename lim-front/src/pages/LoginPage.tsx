import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Box, Typography } from '@mui/material';
import InputField from '../components/InputField';
import Notification from '../components/Notification';
import { AuthService } from '../api/services/authService';


const authService = new AuthService();

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    telephone: '',
    password: '',
  });
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [openSnackbar, setOpenSnackbar] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setMessage(null);

    try {
      await authService.login(formData, setError);
      setMessage('Вход выполнен успешно!');
      setOpenSnackbar(true);
      setTimeout(() => navigate('/schedule'), 1500);
    } catch (err: any) {
      setMessage(null);
      setOpenSnackbar(true);
      setError(err);
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
        Вход
      </Typography>

      <form onSubmit={handleSubmit} style={{ width: '100%', maxWidth: '400px' }}>
        <InputField
          label="Телефон"
          name="telephone"
          type="tel"
          value={formData.telephone}
          onChange={handleChange}
          required
        />
        <InputField
          label="Пароль"
          name="password"
          type="password"
          value={formData.password}
          onChange={handleChange}
          required
        />
        <Box sx={{ display: 'flex', justifyContent: 'space-between', marginBottom: 2 }}>
          <Button type="submit" variant="contained" color="primary" sx={{ width: '48%' }}>
            Войти
          </Button>
          <Button variant="text" color="inherit" onClick={() => navigate('/signup')} sx={{ width: '48%' }}>
            Регистрация
          </Button>
        </Box>
      </form>

      <Notification
        open={openSnackbar}
        message={error || message}
        severity={error ? 'error' : 'success'}
        onClose={() => setOpenSnackbar(false)}
      />
    </Box>
  );
};

export default LoginPage;
