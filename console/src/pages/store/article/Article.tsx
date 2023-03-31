import {useParams} from 'react-router-dom';
import {useEffect, useState} from 'react';
import {Divider} from 'react-windy-ui';
import {get} from '@/client/Request';
import {ArticleInfo} from '@/Types';
import dayjs from 'dayjs';
import {buildUrl} from '@/common/utils';

export default function Article() {
    const {id} = useParams();
    const [article, set] = useState<ArticleInfo | null>(null);

    useEffect(() => {
        get(buildUrl(`articles/${id}`)).then((data) => {
            set(data as ArticleInfo);
        });
    }, [id]);

    return (
        <div className="bs-article">
            <div>
                <h2>{article?.name}</h2>
                <h4 className="text comment">{dayjs(article?.createDate).format('YYYY-MM-DD')}</h4>
                <Divider/>
                <div dangerouslySetInnerHTML={{__html: article?.content ?? ''}}></div>
            </div>
        </div>
    );
}
