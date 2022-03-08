import './Thumbnail.scss';

export function Thumbnail(props: { url: string }) {

    const style: React.CSSProperties = {
        backgroundImage: `url(${props.url})`,
    };

    return (
        <div className="thumbnail" style={style}></div>
    )
}