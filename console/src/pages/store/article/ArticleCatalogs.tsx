import { Navbar } from 'react-windy-ui';
import { useContext } from 'react';
import { ArticleCatalogContext } from '@/common/Context';
import { Catalog } from '@/Types';

export default function ArticleCatalogs() {
  const articleCatalogCtx = useContext(ArticleCatalogContext);
  const { catalogs, changeArticleCatalog, selectedCatalogId } = articleCatalogCtx;

  return (
    <Navbar type="primary" extraClassName="bs-header-bar" hasBorder={false} hasBox={false}>
      <Navbar.Title>书库</Navbar.Title>
      <Navbar.List justify="start">
        {catalogs.map((c: Catalog, index: number) => {
          return (
            <Navbar.Item
              key={`c-${c.id}`}
              hasBackground={true}
              hasBar={true}
              onClick={() => changeArticleCatalog && changeArticleCatalog(c.id)}
              active={c.id === selectedCatalogId}>
              {c.name}
            </Navbar.Item>
          );
        })}
        <Navbar.Item
          hasBackground
          hasBar
          onClick={() => changeArticleCatalog && changeArticleCatalog('all')}
          active={'all' === selectedCatalogId}>
          全部文章
        </Navbar.Item>
      </Navbar.List>
    </Navbar>
  );
}
