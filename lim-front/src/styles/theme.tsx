// theme.ts
import { createTheme } from '@mui/material/styles';

// Устанавливаем цвета
const theme = createTheme({
  palette: {
    primary: {
      main: '#7D7D7D', // цвет для кнопок и ссылок
    },
    secondary: {
        main: '#E0E0E0', // цвет для кнопок и ссылок
      },
    background: {
      default: '#FFFFFF', // фон для страницы
      paper: '#E0E0E0', // фоновые элементы в блоках
    },
    text: {
      primary: '#444444', // основной цвет текста на светлом фоне
      secondary: '#FFFFFF', // текст на темном фоне
    },
  },
  typography: {
    // Шрифт Cormorant
    fontFamily: '"Cormorant", serif', 
    h1: {
      fontFamily: '"Cormorant", serif',
      fontWeight: 700, // Cormorant Bold
      fontSize: '36px',
    },
    h2: {
      fontFamily: '"Cormorant", serif',
      fontWeight: 400, // Cormorant Regular
      fontSize: '28px',
    },
    h3: {
      fontFamily: '"Cormorant", serif',
      fontWeight: 300, // Cormorant Light
      fontSize: '24px',
    },
    h4: {
        fontFamily: '"Cormorant", serif',
        fontWeight: 300, // Cormorant Light
        fontSize: '20px',
    },
    h5: {
        fontFamily: '"Cormorant", serif',
        fontWeight: 300, // Cormorant Light
        fontSize: '18px',
    },

    h6: {
        fontFamily: '"Cormorant", serif',
        fontWeight: 300, // Cormorant Light
        fontSize: '14px',
    },
      
    body1: {
      fontFamily: '"Cormorant", serif', // для текста
      fontWeight: 400, // Cormorant Regular
      fontSize: '16px',
    },
    body2: {
      fontFamily: '"Cormorant", serif',
      fontWeight: 300, // Cormorant Light
      fontSize: '16px',
    },
    button: {
      fontFamily: '"Cormorant", serif',
      fontWeight: 700, // Cormorant Bold
      fontSize: '16px',
    },
  },
});

export default theme;
