import { Col, Row } from 'antd';
import React from 'react';
import './index.css'
import Navigation from "./Navigation/Navigation";
import {Link, Outlet} from "react-router-dom";

// row 下分成24等分的col
export const HomePage: React.FC = () => (
    <div className='homepage'>
        <Row className='header'>
            <Col span={2} className='title'>
                小说管理后台
            </Col>
            <Col span={20}>
            </Col>
            <Col span={1} className='userCenter'>
                <Link to="/loginOut">LoginOut</Link>
            </Col>
            <Col span={1} className='userCenter'>
                用户信息
            </Col>
        </Row>
        <Row className='body'>
            <Col span={3} style={{border: "solid"}}>
                <Navigation></Navigation>
            </Col>
            <Col span={21}>
                <Outlet />
            </Col>
        </Row>
    </div>
);