import React, {Component} from 'react';
import {Input, Menu, Segment} from 'semantic-ui-react'
import 'semantic-ui-css/semantic.min.css';
import './App.css';
import Songs from "./Songs/Songs";

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {categories: [], activeItem: '', activeCategory: []}
    }

    componentDidMount() {
        this.fetchSongs();
    }

    fetchSongs() {
        fetch('http://florianboulay.pro:4000/api/songs')
            .then(res => res.json())
            .then((categories) => {
                this.setState({categories: categories, activeItem: categories[0].name, activeCategory: categories[0]});
            });
    }

    handleItemClick = (e, {name}) => {
        const activeCategory = this.state.categories.find((item) => {
            return item.name === name;
        });
        this.setState({
            activeItem: name,
            activeCategory: activeCategory
        });
    };

    render() {

        return (
            <div>
                <div className={'ui container'}>
                    <Menu pointing>
                        {this.state.categories.map((item) =>
                            <Menu.Item name={item.name} active={this.state.activeItem === item.name}
                                       key={item.name}
                                       onClick={this.handleItemClick}/>
                        )}
                    </Menu>
                    <Songs refresh={this.fetchSongs} category={this.state.activeCategory}/>
                </div>
            </div>
        );
    }
}

export default App;
