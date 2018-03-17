import React, {Component} from 'react';
import {Button} from "semantic-ui-react";

export default class Songs extends Component {

    constructor(props) {
        super(props);
    }

    playSong(event) {
        fetch('http://localhost:4000/api/songs?category=' + this.props.category.name + '&song='+event.target.value);
    }

    render() {
        if (!this.props.category.songs) {
            return null;
        }
        return (
            <div className={'ui segments'}>
                <div className={'ui grid container'}>
                    <div className={'ui seven column grid'}>
                        {this.props.category.songs.map((song) =>
                            <div className={'column'}>
                                <Button key={song.name} value={song.name}
                                        onClick={(event) => this.playSong(event)}>{song.name}</Button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        );
    }
}