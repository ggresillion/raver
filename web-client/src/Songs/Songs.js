import React, { Component } from 'react';
import { Button } from "semantic-ui-react";

export default class Songs extends Component {

    constructor(props) {
        super(props);
        this.state = { songsToUpload: [] };
        this.onChange = this.onChange.bind(this);
        this.uploadSongs = this.uploadSongs.bind(this);
    }

    playSong(event) {
        fetch('http://florianboulay.pro:4000/api/songs?category=' + this.props.category.name + '&song=' + event.target.value);
    }

    uploadSongs(e) {
        console.log('uploading ...');
        e.preventDefault();
        let data = new FormData();
        data.append('category', this.props.category.name);
        data.append('songs', this.state.songsToUpload);
        fetch('http://florianboulay.pro:4000/api/songs', {
            method: 'POST',
            body: data
        })
            .then((res) => {
                //this.props.refresh();
            })
    }

    onChange(e) {
        this.setState({ songsToUpload: e.target.files[0] })
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
                            <form onSubmit={this.uploadSongs}>
                                <input name="songs" type="file" onChange={this.onChange} style={{ display: 'block' }} />
                                <Button basic color='green' type='submit' key='add'>+</Button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}