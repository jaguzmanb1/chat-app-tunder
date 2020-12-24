import React, { Component } from 'react';

export default class Chat extends Component {
    constructor(props){
        super(props);
        this.state = {
            textMessage: '',
            to: 0,
            from: 0,
            connected: false
        }

        this.handleTextMessageChange = this.handleTextMessageChange.bind(this);
        this.handleToChange = this.handleToChange.bind(this);
        this.handleFromChange = this.handleFromChange.bind(this);
        this.makeConnect = this.makeConnect.bind(this);
        this.sendMessage = this.sendMessage.bind(this);

    }
    ws = new WebSocket('ws://localhost')

    handleTextMessageChange(event) {
        this.setState({
            textMessage: event.target.value
        })
    }

    handleToChange(event) {
        this.setState({
            to: parseInt(event.target.value)
        })
    }

    handleFromChange(event) {
        this.setState({
            from: parseInt(event.target.value)
        })
    }

    sendMessage(){
        this.ws.send(JSON.stringify({
            "to": this.state.to,
            "from": this.state.from,
            "message": this.state.textMessage
        }))
    }

    makeConnect () {
        this.ws = new WebSocket('ws://localhost:3030/ws/' + this.state.from)

        this.ws.onopen = () => {
            // on connecting, do nothing but log it to the console
            console.log('connected')
            this.setState( {
                connected: true
            })
        }

        this.ws.onmessage = evt => {
            // listen to data sent from the websocket server
            const message = JSON.parse(evt.data)
            this.setState({dataFromServer: message})
            console.log(message)
        }

        this.ws.onclose = () => {
            console.log('disconnected')
            this.setState( {
                connected: false
            })
            // automatically try to reconnect on connection loss
        }

    }

    render () {
        return (
            (this.state.connected ? 

            <div>
                <p>Chat works</p>
                To 
                <input type="text" value={this.state.to} onChange={this.handleToChange} />

                Message 
                <input type="text" value={this.state.textMessage} onChange={this.handleTextMessageChange} />

                <button onClick={this.sendMessage}> Enviar </button>
            </div> 
            : 
            <div>
                From 
                <input type="text" value={this.state.from} onChange={this.handleFromChange} />
                <button onClick={this.makeConnect}> Connect </button>

            </div>
            )
            
        )
    }


}