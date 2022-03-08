import './Loader.scss';

export function Loader() {

    return (
        <div className='loader'>
            <div className="lds-ring"><div></div><div></div><div></div><div></div></div>
        </div>
    );
}