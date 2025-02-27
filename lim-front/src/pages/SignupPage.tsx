import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthService } from '../api/services/authService';
import { Button, Box, Typography } from '@mui/material';
import Notification from '../components/Notification';
import InputField from '../components/InputField';

const RegistrationPage: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    telephone: '',
    email: '',
    password: '',
  });
  const [message, setMessage] = useState<string | null>(null);
  const [openSnackbar, setOpenSnackbar] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const authService = new AuthService();

  
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      await authService.signup(formData, setError);
      setMessage('Регистрация прошла успешно!');
      setError(null)
      setOpenSnackbar(true);
      setTimeout(() => navigate('/login'), 1500);
    } catch (error: any) {
      setError(error);
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
        Регистрация
      </Typography>

      <form onSubmit={handleSubmit} style={{ width: '100%', maxWidth: '400px' }}>
        <InputField
          label="Имя"
          name="name"
          type="text"
          value={formData.name}
          onChange={handleChange}
          required
        />
        <InputField
          label="Телефон"
          name="telephone"
          type="tel"
          value={formData.telephone}
          onChange={handleChange}
          required
        />
        <InputField
          label="Email"
          name="email"
          type="email"
          value={formData.email}
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
          <Button
            type="submit"
            variant="contained"
            color="primary"
            sx={{ width: '48%' }}
          >
            Зарегистрироваться
          </Button>
          <Button
            variant="text"
            color="inherit"
            onClick={() => navigate('/login')}
            sx={{ width: '48%' }}
          >
            Уже есть аккаунт?
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

export default RegistrationPage;
