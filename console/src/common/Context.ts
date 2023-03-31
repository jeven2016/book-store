import React from 'react';
import { ArticleCatalogContextValue, WindowContextValue } from '@/Types';

export const WindowContext = React.createContext<WindowContextValue>({ mdWindow: false });

export const ArticleCatalogContext = React.createContext<ArticleCatalogContextValue>({
  catalogs: [],
  selectedCatalogId: null,
  changeArticleCatalog: null
});
