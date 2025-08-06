# Student Attendance Web Application

A modern, responsive web login interface for the Student Attendance API built with React, TypeScript, and Tailwind CSS.

## Features

### 🔐 **Authentication**
- **Multi-user login**: Support for admin, teacher, and student login types
- **JWT token management**: Automatic token handling and storage
- **Session persistence**: Login state persists across browser sessions
- **Cookie support**: Supports both JWT tokens and HTTP-only cookies

### 🌍 **Internationalization (i18n)**
- **Multi-language support**: English and Bahasa Indonesia
- **Dynamic language switching**: Change language without page reload
- **Browser language detection**: Automatically detects user's preferred language
- **Persistent language preference**: Saves language choice in localStorage

### 🌙 **Dark Mode**
- **Theme options**: Light, Dark, and System themes
- **Smooth transitions**: Animated theme switching
- **System preference detection**: Follows OS dark mode settings
- **Persistent theme choice**: Saves theme preference in localStorage

### 🎨 **Modern UI/UX**
- **Responsive design**: Works on desktop, tablet, and mobile
- **Tailwind CSS**: Modern, utility-first CSS framework
- **Smooth animations**: Fade-in and slide-up animations
- **Form validation**: Real-time validation with error messages
- **Loading states**: Visual feedback during API calls

### 🛠 **Technical Features**
- **TypeScript**: Full type safety and better development experience
- **Zustand**: Lightweight state management
- **React Hook Form**: Performant form handling with validation
- **Fetch API**: Modern HTTP client with comprehensive error handling
- **Try-catch error handling**: Robust error management throughout

## Technology Stack

- **React 19**: Latest React with TypeScript
- **Zustand**: State management
- **Tailwind CSS**: Styling framework
- **React Hook Form**: Form handling
- **i18next**: Internationalization
- **Lucide React**: Modern icon library
- **Vite**: Fast build tool

## Installation

1. **Navigate to web directory**:
   ```bash
   cd web
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Start development server**:
   ```bash
   npm run dev
   ```

4. **Open in browser**:
   Visit `http://localhost:5173`

## Configuration

### API Endpoint
Update the API base URL in `src/services/api.ts`:
```typescript
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

### Supported Languages
Add new languages in `src/i18n/locales/` and update the configuration in `src/i18n/index.ts`.

## Usage

### Login Process
1. **Select user type**: Admin, Teacher, or Student
2. **Enter credentials**:
   - **Admin**: Email address
   - **Teacher**: Teacher ID
   - **Student**: Student ID
3. **Enter password**
4. **Submit form**

### Theme Switching
- Click the theme toggle button in the top-right corner
- Choose from Light, Dark, or System theme

### Language Switching
- Click the language selector in the top-right corner
- Choose from available languages (English/Bahasa Indonesia)

## API Integration

The application integrates with the Student Attendance API:

### Endpoints Used
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/auth/logout` - User logout

### Error Handling
- **Network errors**: Connection failures
- **Authentication errors**: Invalid credentials
- **Server errors**: API server issues
- **Validation errors**: Form validation failures

### Authentication Flow
1. User submits login form
2. API call made to `/auth/login`
3. JWT token received and stored
4. Token cached in localStorage and set as HTTP-only cookie
5. Subsequent API calls include Authorization header
6. Automatic logout on token expiration

## Project Structure

```
src/
├── components/          # React components
│   ├── LoginForm.tsx   # Main login form
│   ├── LoginPage.tsx   # Complete login page
│   ├── LanguageSwitcher.tsx
│   └── ThemeToggle.tsx
├── stores/             # Zustand stores
│   ├── authStore.ts    # Authentication state
│   └── themeStore.ts   # Theme state
├── services/           # API services
│   └── api.ts          # API client with error handling
├── types/              # TypeScript type definitions
│   └── auth.ts         # Auth-related types
├── i18n/               # Internationalization
│   ├── index.ts        # i18n configuration
│   └── locales/        # Translation files
│       ├── en.json     # English translations
│       └── id.json     # Indonesian translations
├── App.tsx             # Main application component
├── main.tsx            # Application entry point
└── index.css           # Global styles with Tailwind
```

## Development

### Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### State Management
The application uses Zustand for state management with persistence:

- **Auth Store**: User authentication state, login/logout functions
- **Theme Store**: Theme preference and dark mode toggle

### Styling
- **Tailwind CSS**: Utility-first CSS framework
- **Custom CSS**: Additional animations and components
- **Dark mode**: Class-based dark mode implementation
- **Responsive**: Mobile-first responsive design

## Security Features

- **Token management**: Secure JWT token handling
- **HTTP-only cookies**: Additional security layer
- **Input validation**: Client-side form validation
- **Error handling**: Secure error message display
- **Session management**: Automatic token cleanup

## Browser Support

- **Modern browsers**: Chrome, Firefox, Safari, Edge
- **Responsive design**: Works on all screen sizes
- **Progressive enhancement**: Graceful degradation for older browsers

## Contributing

1. Follow TypeScript best practices
2. Use Tailwind CSS classes for styling
3. Add translations for new text content
4. Test in both light and dark modes
5. Ensure responsive design on all screen sizes