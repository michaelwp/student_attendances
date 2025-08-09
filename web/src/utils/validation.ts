export interface PasswordValidationRule {
  test: (password: string) => boolean;
  message: string;
}

export const passwordValidationRules: PasswordValidationRule[] = [
  {
    test: (password: string) => password.length >= 9,
    message: 'Password must be at least 9 characters long',
  },
  {
    test: (password: string) => /[A-Z]/.test(password),
    message: 'Password must contain at least one uppercase letter',
  },
  {
    test: (password: string) => /[a-z]/.test(password),
    message: 'Password must contain at least one lowercase letter',
  },
  {
    test: (password: string) => /[0-9]/.test(password),
    message: 'Password must contain at least one number',
  },
  {
    test: (password: string) => /[!@#$%^&*(),.?":{}|<>]/.test(password),
    message: 'Password must contain at least one symbol (!@#$%^&*(),.?":{}|<>)',
  },
];

export const validatePassword = (password: string): string[] => {
  const errors: string[] = [];
  passwordValidationRules.forEach(rule => {
    if (!rule.test(password)) {
      errors.push(rule.message);
    }
  });
  return errors;
};

export const isPasswordValid = (password: string): boolean => {
  return passwordValidationRules.every(rule => rule.test(password));
};

export const getPasswordStrength = (password: string): number => {
  const passedRules = passwordValidationRules.filter(rule => rule.test(password)).length;
  return (passedRules / passwordValidationRules.length) * 100;
};

// Email validation
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

// Phone number validation (basic)
export const isValidPhone = (phone: string): boolean => {
  const phoneRegex = /^\+?[\d\s\-()]{10,}$/;
  return phoneRegex.test(phone);
};

// Student/Teacher ID validation
export const isValidUserID = (id: string): boolean => {
  // Only alphanumeric characters, 3-20 characters long
  const userIdRegex = /^[A-Za-z0-9]{3,20}$/;
  return userIdRegex.test(id);
};