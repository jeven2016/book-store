import {
  Button,
  Col,
  Container,
  Divider,
  IconSearch,
  IconStarBorder,
  Input,
  InputGroup,
  Responsive,
  Row,
  Space,
  useMediaQuery
} from 'react-windy-ui';
import ArticleCatalogs from './ArticleCatalogs';
import Content from '@/pages/store/article/Content';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { ArticleCatalogContext, ArticleSearchCtx } from '@/common/Context';
import { Catalog } from '@/Types';
import { buildUrl } from '@/common/utils';
import { get } from '@/client/Request';

export default function Home() {
  const [catalogs, setCatalogs] = useState<Catalog[]>([]);
  const [selectedCatalogId, setSelectedCatalogId] = useState<string | null>(null);
  const [search, setSearch] = useState<string>('');
  const [searchText, setSearchText] = useState<string>('');

  useEffect(() => {
    get(buildUrl('/article-catalogs')).then((data) => {
      const list = data as Catalog[];
      setCatalogs(list);
    });
  }, []);

  useEffect(() => {
    if (!selectedCatalogId && catalogs.length > 0) {
      setSelectedCatalogId(catalogs[0].id);
    }
  }, [catalogs, selectedCatalogId]);

  //切换catalog
  const changeArticleCatalog = useCallback((id: string) => setSelectedCatalogId(id), []);
  const providerValue = useMemo(
    () => ({
      catalogs,
      selectedCatalogId,
      changeArticleCatalog
    }),
    [catalogs, selectedCatalogId]
  );

  const searchArticles = useCallback(() => {}, []);

  return (
    <ArticleCatalogContext.Provider value={providerValue}>
      <Header
        search={search}
        setSearch={setSearch}
        changeArticleCatalog={changeArticleCatalog}
        setSearchText={setSearchText}
      />
      <ArticleSearchCtx.Provider value={{ name: searchText }}>
        <Content />
      </ArticleSearchCtx.Provider>
    </ArticleCatalogContext.Provider>
  );
}

function Header(props) {
  const { search, setSearch, changeArticleCatalog, setSearchText } = props;
  const { matches } = useMediaQuery(Responsive.sm.max);

  const doSearch = useCallback(() => {
    changeArticleCatalog('all');
    setSearchText(search);
  }, [search]);

  return (
    <>
      <div className="bs-header-wrapper">
        <Container extraClassName="bs-header-wrapper" autoAdjust={true}>
          <Row justify="center" align="center" extraClassName="bs-header">
            {!matches && (
              <Col extraClassName="bs-header-logo" align="center" flexCol={true} col={2}>
                图书馆
              </Col>
            )}
            <Col col={matches ? 12 : 8} justify="center" flexCol={true}>
              <InputGroup size="small" extraClassName="bs-header-search-input">
                <Input
                  extraClassName="bs-input-search"
                  placeholder="书名、作者"
                  value={search}
                  onKeyDown={(e) => e.key === 'Enter' && doSearch()}
                  onChange={(e) => setSearch(e.target.value)}
                />
                <InputGroup.Item autoScale={false}>
                  <Button
                    extraClassName="bs-btn-search"
                    onClick={() => {
                      doSearch();
                    }}>
                    <IconSearch />
                  </Button>
                </InputGroup.Item>
              </InputGroup>
            </Col>
            <Col co={matches ? 12 : 2} style={{ marginTop: matches ? '.5rem' : '0' }}>
              <Space>
                <Button leftIcon={<IconStarBorder />} hasBorder={false} hasBox={false} size="small">
                  我的书架
                </Button>
                <Divider direction="vertical" />
                <Button hasBox={false} hasBorder={false} size="small">
                  登录
                </Button>
              </Space>
            </Col>
          </Row>
        </Container>
      </div>
      <div className="bs-navbar">
        <Container autoAdjust={true}>
          <ArticleCatalogs />
        </Container>
      </div>
    </>
  );
}
