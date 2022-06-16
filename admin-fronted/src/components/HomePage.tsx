import { Col, Row } from 'antd';
import React from 'react';

const HomePage: React.FC = () => (
    <>
        <Row style={{ background: 'skyblue' }}>
            <Col span={1} style={{ fontSize: '40px', }} >Admin</Col>
        </Row>
    </>
);

export default HomePage;