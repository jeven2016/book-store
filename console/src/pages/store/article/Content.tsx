import { Box, Container, Menu, Responsive, useMediaQuery } from 'react-windy-ui';
import ArticleList from '@/pages/store/article/ArticleList';
import { useCallback, useContext, useEffect, useState } from 'react';
import { ArticleCatalogContext, WindowChangeContext } from '@/common/Context';
import { ArticleCatalogContextValue, Catalog, WindowChangeInfo } from '@/Types';
import AllArticles from '@/pages/store/article/AllArticles';

export default function Content() {
  const [selectedSubCatalogId, setSelectedSubCatalogId]: [
    string | null,
    (value: ((prevState: string | null) => string | null) | string | null) => void
  ] = useState<string | null>(null);
  const articleCatalogCtx: ArticleCatalogContextValue = useContext(ArticleCatalogContext);
  const { catalogs, selectedCatalogId }: ArticleCatalogContextValue = articleCatalogCtx;
  const { sm: matches }: WindowChangeInfo = useContext<WindowChangeInfo>(WindowChangeContext);

  let chd = [] as Catalog[];

  if (catalogs.length > 0 && selectedCatalogId) {
    const selectedCatalogs = catalogs.filter((c) => c.id === selectedCatalogId);
    if (selectedCatalogs.length > 0) {
      chd = selectedCatalogs[0].children;
    }
  }

  useEffect(() => {
    if (selectedSubCatalogId && chd.some((c) => c.id === selectedSubCatalogId)) {
      setSelectedSubCatalogId(selectedSubCatalogId);
      return;
    }

    if (chd.length > 0) {
      setSelectedSubCatalogId(chd[0].id);
    }
  }, [selectedSubCatalogId, chd]);

  const clickItem = useCallback((selectedId) => {
    setSelectedSubCatalogId(selectedId);
  }, []);

  const menu = (
    <div className="bs-content-inner-left">
      <Menu
        direction={matches ? 'horizontal' : 'vertical'}
        type="primary"
        hasBox={false}
        primaryBarPosition="left"
        activeItems={selectedSubCatalogId ?? ''}
        onClickItem={clickItem}>
        {chd.map((c: Catalog) => (
          <Menu.Item key={c.id} id={c.id}>
            {c.name}
          </Menu.Item>
        ))}
      </Menu>
    </div>
  );

  return selectedCatalogId === 'all' ? (
    <AllArticles />
  ) : (
    <Container extraClassName="bs-content" autoAdjust={true}>
      <div className="bs-content-panel">
        <div className="bs-panel-title">
          <h3>精品分类</h3>
        </div>
        {matches && (
          <div>
            {menu}
            <ArticleList smWindow={matches} subCatalogId={selectedSubCatalogId} />
          </div>
        )}
        {!matches && (
          <Box
            extraClassName="bs-content-inner"
            block
            align="start"
            left={menu}
            center={<ArticleList smWindow={matches} subCatalogId={selectedSubCatalogId} />}
          />
        )}
      </div>
    </Container>
  );
}
