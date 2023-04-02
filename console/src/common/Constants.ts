export const URI = {
  defaultHome: process.env.REACT_APP_DEFAULT_HOME ?? '',
  redirectFlag: 'refer=',

  loginPath: '/login'
};

//keycloak auth相关参数
export const AUTH_CONFIG = {
  authMethod: 'check-sso',
  pkceMethod: 'S256'
};

/**
 * 当屏幕窗口变化后，组件的折叠/展开的参数
 */
export const LAYOUT_SETTING = {
  CONTENT: { MIN: '1rem', MAX: '17.625rem' },
  LEFT_SLIDER: { MIN: '-16.625rem', MAX: '0rem' }
};

export const BACKEND_API_PREFIX = '/api/v1';
