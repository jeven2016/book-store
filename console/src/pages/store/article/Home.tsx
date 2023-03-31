import {
    Button,
    Col,
    Container,
    Divider,
    IconSearch,
    IconStarBorder,
    Input,
    InputGroup,
    Row,
    Space
} from 'react-windy-ui';
import ArticleCatalogs from './ArticleCatalogs';
import Content from '@/pages/store/article/Content';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {ArticleCatalogContext} from '@/common/Context';
import {Catalog} from '@/Types';
import {buildUrl} from '@/common/utils';
import {get} from '@/client/Request';

export default function Home() {
    const [catalogs, setCatalogs] = useState<Catalog[]>([]);
    const [selectedCatalogId, setSelectedCatalogId] = useState<string | null>(null);

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

    return (
        <ArticleCatalogContext.Provider value={providerValue}>
            <Header/>
            <Content/>
        </ArticleCatalogContext.Provider>
    );
}

function Header() {
    return (
        <>
            <div className="bs-header-wrapper">
                <Container extraClassName="bs-header-wrapper" autoAdjust={true}>
                    <Row justify="center" align="center" extraClassName="bs-header">
                        <Col extraClassName="bs-header-logo" align="center" flexCol={true} col={2}>
                            图书馆
                        </Col>
                        <Col col={8} justify="center" flexCol={true}>
                            <InputGroup size="small" extraClassName="bs-header-search-input">
                                <Input extraClassName="bs-input-search" placeholder="书名、作者"/>
                                <InputGroup.Item autoScale={false}>
                                    <Button extraClassName="bs-btn-search">
                                        <IconSearch/>
                                    </Button>
                                </InputGroup.Item>
                            </InputGroup>
                        </Col>
                        <Col co={2}>
                            <Space>
                                <Button leftIcon={<IconStarBorder/>} hasBorder={false} hasBox={false} size="small">
                                    我的书架
                                </Button>
                                <Divider direction="vertical"/>
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
                    <ArticleCatalogs/>
                </Container>
            </div>
        </>
    );
}
