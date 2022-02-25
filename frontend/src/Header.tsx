import React from "react";
import './Header.css'

export class Header extends React.Component<any, { res: string }> {

    async componentDidMount() {
        const res = await (await fetch('http://localhost:8080/')).text();
        this.setState({ res });
    }

    render() {
        const res = this.state?.res;
        if (!res) {
            return (
                <div className="header">Loading...</div>
            )
        }
        return (
            <div className="header">
                {res}
            </div>
        )
    }
}