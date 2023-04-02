import { BACKEND_API_PREFIX } from '@/common/Constants';

export const buildUrl = (uri = '') => {
  if (uri.length == 0) {
    return BACKEND_API_PREFIX;
  }
  let suffix = uri;
  if (!uri.startsWith('/')) {
    suffix = `/` + uri;
  }
  return `${BACKEND_API_PREFIX}${suffix}`;
};
