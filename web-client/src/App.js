import React, {Component} from 'react';
import {Input, Menu, Segment, Button} from 'semantic-ui-react'
import 'semantic-ui-css/semantic.min.css';
import './App.css';
import Songs from "./Songs/Songs";
import config from './config'

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {categories: [], activeItem: '', activeCategory: []}
        this.fetchSongs = this.fetchSongs.bind(this);
    }

    componentDidMount() {
        this.fetchSongs();
    }

    fetchSongs() {
        fetch(config.api.url + '/api/songs')
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
                        {/*<Menu.Item> <Button basic color="green" name="+" key="+" onClick={console.log}>+</Button></Menu.Item>*/}
                    </Menu>
                    <Songs refresh={this.fetchSongs} category={this.state.activeCategory}/>
                </div>
            </div>
        );
    }
}

export default App;
