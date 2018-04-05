import React, { Component } from 'react'
import { Button, Header, Icon, Modal, Input } from 'semantic-ui-react'

export default class ModalExampleControlled extends Component {
    state = { modalOpen: false }

    handleOpen = () => this.setState({ modalOpen: true })

    handleClose = () => this.setState({ modalOpen: false })

    onChange= (e)=> this.setState({categoryName: e.target.value})

    handleSend = () => {
        console.log(this.state.categoryName);
        this.handleClose();
    }

    render() {
        const inlineStyle = {
            modal: {
                marginTop: '0px !important',
                marginLeft: 'auto',
                marginRight: 'auto'
            }
        };
        return (
            <Modal
                trigger={<Button basic color="green" name="+" key="+" onClick={this.handleOpen}>+</Button>}
                open={this.state.modalOpen}
                onClose={this.handleClose}
                style={inlineStyle.modal}
                size='tiny'
            >
                <Header icon='browser' content='New Category' />
                <Modal.Content>
                    <h3>Enter a name for the new category :</h3>
                    <Input onChange={this.onChange} placeholder='Category name' />
                </Modal.Content>
                <Modal.Actions>
                    <Button negative onClick={this.handleClose} content="Cancel"/>
                    <Button positive onClick={this.handleSend} icon='checkmark' labelPosition='right' content='Create' />
                </Modal.Actions>
            </Modal>
        )
    }
}