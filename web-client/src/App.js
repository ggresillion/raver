import React, {Component} from 'react';
import {Button} from 'semantic-ui-react';
import logo from './logo.svg';
import './App.css';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {songs: []}
    }

    componentDidMount() {
        this.fetchSongs();
    }

    fetchSongs() {
        fetch('http://localhost:4000/api/songs')
            .then(res => res.json())
            .then((songs) => {
                this.setState({songs: songs});
            });
    }

    playSong(event){
        fetch('http://localhost:4000/api/songs/'+event.target.value);
    }

    render() {

        return (
            <div>
                {this.state.songs.map((song) => <Button key={song.name} value={song.name} onClick={this.playSong}>{song.name}</Button>)}
            </div>
        );
    }
}

export default App;
