import React, { Component } from 'react';
import { Button, Grid, GridColumn, Segment } from "semantic-ui-react";
import config from "../config";
import './Songs.css';

export default class Songs extends Component {

    constructor(props) {
        super(props);
        this.state = { songsToUpload: [] };
        this.uploadSongs = this.uploadSongs.bind(this);
        this.openUploadDialog = this.openUploadDialog.bind(this);
    }

    playSong(event) {
        fetch(config.api.url + '/api/songs?category=' + this.props.category.name + '&song=' + event.target.value);
    }

    uploadSongs(e) {
        console.log('uploading ...');
        e.preventDefault();
        let data = new FormData();
        data.append('category', this.props.category.name);
        data.append('songs', e.target.files[0]);
        fetch(config.api.url + '/api/songs', {
            method: 'POST',
            body: data
        })
            .then((res) => {
                this.props.refresh();
            })
    }

    openUploadDialog() {
        this.refs.input.click();
    }

    render() {
        if (!this.props.category.songs) {
            return null;
        }
        return (
            <Segment id="song-grid">
                        {this.props.category.songs.map((song) =>
                                <Button id="song-button" key={song.name} value={song.name}
                                    onClick={(event) => this.playSong(event)}>{song.name}</Button>
                        )}
                            <input ref="input" name="songs" type="file" onChange={this.uploadSongs} style={{ display: 'none' }} />
                            <Button id="song-button" basic color='green' onClick={this.openUploadDialog} key='add'>+</Button>
            </Segment>
        );
    }
}