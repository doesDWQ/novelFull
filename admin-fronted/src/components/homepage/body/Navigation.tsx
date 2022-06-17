import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import React, { useState } from 'react';
import {
    useNavigate
} from "react-router-dom";

type MenuItem = Required<MenuProps>['items'][number];

function getItem(
    label: React.ReactNode,
    key?: React.Key | null,
    icon?: React.ReactNode,
    children?: MenuItem[],
    type?: 'group',
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
        type,
    } as MenuItem;
}

const items: MenuItem[] = [
    getItem('首页','/index', null),
    getItem('用户管理', '/user', null, [
        getItem('管理员管理', '/users'),
        getItem('普通用户管理', '/adminUsers'),
    ]),

    // getItem('Navigation Two', 'sub2', <AppstoreOutlined />, [
    //     getItem('Option 5', '5'),
    //     getItem('Option 6', '6'),
    //     getItem('Submenu', 'sub3', null, [getItem('Option 7', '7'), getItem('Option 8', '8')]),
    // ]),
    //
    // getItem('Navigation Three', 'sub4', <SettingOutlined />, [
    //     getItem('Option 9', '9'),
    //     getItem('Option 10', '10'),
    //     getItem('Option 11', '11'),
    //     getItem('Option 12', '12'),
    // ]),
];

const Navigation: React.FC = () => {
    const [current, setCurrent] = useState('1');
    // 设置主题颜色, 只有黑色和白色
    const theme = 'dark'
    let navigate = useNavigate();

    const onClick: MenuProps['onClick'] = e => {
        console.log('click ', e);
        setCurrent(e.key);
        navigate(e.key)
    };

    return (
        <>
            <Menu style={{height: "100%"}}
                theme={theme}
                onClick={onClick}
                defaultOpenKeys={['sub1']}
                selectedKeys={[current]}
                mode="inline"
                items={items}
            />
        </>
    );
};

export default Navigation;