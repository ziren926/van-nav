import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Button } from "antd";
import { ArrowLeftOutlined, LinkOutlined } from "@ant-design/icons";
import { getLogoUrl } from "../../utils/check";
import { getJumpTarget } from "../../utils/setting";
import "./index.css";

const ToolDetail = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { title, url, des, logo, catelog } = location.state || {};

  // 如果没有数据，返回首页
  if (!title) {
    navigate('/');
    return null;
  }

  return (
    <div className="tool-detail">
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate(-1)}
        className="back-button"
      >
        返回
      </Button>

      <div className="tool-detail-content">
        <div className="tool-detail-header">
          <img
            src={url === "admin" ? logo : getLogoUrl(logo)}
            alt={title}
            className="tool-detail-logo"
          />
          <div className="tool-detail-title">
            <h1>{title}</h1>
            <span className="tool-detail-category">{catelog}</span>
          </div>
        </div>

        <div className="tool-detail-description">
          <h2>描述</h2>
          <p>{des}</p>
        </div>

        <div className="tool-detail-actions">
          <Button
            type="primary"
            icon={<LinkOutlined />}
            href={url}
            target={getJumpTarget() === "blank" ? "_blank" : "_self"}
            size="large"
          >
            访问网站
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ToolDetail;