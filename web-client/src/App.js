import React, { Component } from 'react';
import { Input, Menu, Segment, Button } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css';
import './App.css';
import Songs from "./Songs/Songs";
import config from './config';
import NewCategory from './Dialogs/NewCategory';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = { categories: [], activeItem: '', activeCategory: [], editMode: false }
        this.fetchSongs = this.fetchSongs.bind(this);
        this.toggleEditMode = this.toggleEditMode.bind(this);
    }

    componentDidMount() {
        this.fetchSongs();
    }

    fetchSongs() {
        fetch(config.api.url + '/api/songs')
            .then(res => res.json())
            .then((categories) => {
                this.setState({ categories: categories, activeItem: categories[0].name, activeCategory: categories[0] });
            });
    }

    handleItemClick = (e, { name }) => {
        const activeCategory = this.state.categories.find((item) => {
            return item.name === name;
        });
        this.setState({
            activeItem: name,
            activeCategory: activeCategory
        });
    };

    toggleEditMode() {
        this.setState({ editMode: !this.state.editMode });
    }

    render() {

        return (
            <div>
                    <Menu id="menu" pointing>
                        {this.state.categories.map((item) =>
                            <Menu.Item name={item.name} active={this.state.activeItem === item.name}
                                key={item.name}
                                onClick={this.handleItemClick} />
                        )}
                        <Menu.Item><NewCategory /></Menu.Item>
                        <Menu.Item position="right"><Button color="grey" onClick={this.toggleEditMode}>{this.state.editMode ? "Edit done" : "Edit mode"}</Button></Menu.Item>
                    </Menu>
                    <Songs refresh={this.fetchSongs} category={this.state.activeCategory} />
                </div>
        );
    }
}

export default App;
