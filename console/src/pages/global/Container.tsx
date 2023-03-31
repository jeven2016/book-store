import {Outlet} from 'react-router-dom';
import {Loader} from 'react-windy-ui';
import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';

export default function Container() {
    const {t} = useTranslation();
    const [active, setActive] = useState<boolean>(true);
    useEffect(() => {
        let time;
        if (active) {
            time = setTimeout(() => setActive(false), 1000);
        }

        return () => {
            time && clearTimeout(time);
        };
    }, [active]);

    return (
        <>
            <Loader
                type="primary"
                global
                size="small"
                color="white"
                hasDefaultWidth={false}
                modalStyle={{background: '#000'}}
                direction="horizontal"
                active={active}
                text={t('global.page.loading')}
            />
            <Outlet/>
        </>
    );
}
