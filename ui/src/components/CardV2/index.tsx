import { useMemo } from "react";
import "./index.css";
import { getLogoUrl } from "../../utils/check";
import { getJumpTarget } from "../../utils/setting";

// 添加接口定义
interface CardProps {
  id?: string;          // 添加 id 属性
  title: string;
  url: string;
  des: string;
  logo: string;
  catelog: string;
  onClick: () => void;
  index: number;
  isSearching: boolean;
}

const Card: React.FC<CardProps> = ({
  id,                   // 添加 id 参数
  title,
  url,
  des,
  logo,
  catelog,
  onClick,
  index,
  isSearching
}) => {
  const el = useMemo(() => {
    if (url === "admin") {
      return <img src={logo} alt={title} />
    } else {
      return <img src={getLogoUrl(logo)} alt={title} />
    }
  }, [logo, title, url]);

  const showNumIndex = index < 10 && isSearching;

  return (
    <a
      href={url === "toggleJumpTarget" ? undefined : url}
      onClick={(e) => {
        if (id) {          // 添加 id 判断
          e.preventDefault();
          window.location.href = `/tool/${id}`;
        }
        onClick();
      }}
      target={getJumpTarget() === "blank" ? "_blank" : "_self"}
      rel="noreferrer"
      className="card-box"
    >
      {showNumIndex && <span className="card-index">{index + 1}</span>}
      <div className="card-content">
        <div className="card-left">
          {el}
        </div>
        <div className="card-right">
          <div className="card-right-top">
            <span className="card-right-title" title={title}>{title}</span>
            <span className="card-tag" title={catelog}>{catelog}</span>
          </div>
          <div className="card-right-bottom" title={des}>{des}</div>
        </div>
      </div>
    </a>
  );
};

export default Card;