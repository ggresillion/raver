import React, { Component } from 'react';
import { Button } from "semantic-ui-react";
import config from "../config";

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

    openUploadDialog(){
        this.refs.input.click();
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
                        <div className='column'>
                                <input ref="input" name="songs" type="file" onChange={this.uploadSongs} style={{ display: 'none' }} />
                                <Button basic color='green' onClick={this.openUploadDialog} key='add'>+</Button>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}