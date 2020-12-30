import React, { Component } from 'react';
import 'react-chat-elements/dist/main.css';
// MessageBox component
import { MessageBox } from 'react-chat-elements';
import { ChatItem } from 'react-chat-elements'
import { MessageList } from 'react-chat-elements'


export default class Chat extends Component {
    constructor(props){
        super(props);
        this.state = {
            textMessage: '',
            to: '',
            from: '',
            connected: false,
            messages: {},
            currentUser: ''
        }

        this.handleTextMessageChange = this.handleTextMessageChange.bind(this);
        this.handleToChange = this.handleToChange.bind(this);
        this.handleFromChange = this.handleFromChange.bind(this);
        this.makeConnect = this.makeConnect.bind(this);
        this.sendMessage = this.sendMessage.bind(this);

    }

    handleTextMessageChange(event) {
        this.setState({
            textMessage: event.target.value
        })
    }

    handleToChange(event) {
        this.setState({
            to: event.target.value
        })
    }

    handleFromChange(event) {
        this.setState({
            from: event.target.value,
        })
    }

    sendMessage(){
        this.ws.send(JSON.stringify({
            "to": this.state.to,
            "from": localStorage.setItem('phone', this.state.phoneNumber),
            "message": this.state.textMessage
        }))
    }

    componentDidMount(){
        this.makeConnect();
    }

    componentWillUnmount(){
    }

    makeConnect () {
        this.ws = new WebSocket('wss://localhost:3030/ws')

        this.ws.onopen = () => {
            // on connecting, do nothing but log it to the console
            this.ws.send(JSON.stringify({
                "message": localStorage.getItem('token'),
                connected: true
            }))
        }

        this.ws.onmessage = evt => {
            // listen to data sent from the websocket server
            var messages = JSON.parse(evt.data)
            var users = {}

            messages.map(item => {
                var userData = {
                    message: item.message,
                    date: item.date
                }
                if (users[item.from] != null) {
                    users[item.from] = users[item.from].concat(userData) 
                }else {
                    users[item.from] = []
                    users[item.from].concat(userData)
                }
            })

            this.setState({
                messages: users
            })
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
        console.log(this.state.messages)
        return (

            <div>
                <div className="row">
                    <div className="col-4" 
                    style={{
                        paddingRight: 0
                    }} >
                    {
                        Object.keys(this.state.messages).map((key, index) =>{
                            return(<ChatItem
                                avatar={'https://upload.wikimedia.org/wikipedia/commons/7/7c/User_font_awesome.svg'}
                                alt={'Reactjs'}
                                title={key}
                                subtitle={this.state.messages[key][0].message}
                                date={new Date(this.state.messages[key][0].date)}
                                unread={0} 
                            />)
                        })
                    }
                    </div>
                    <div className="col-8" 
                    style={{
                        paddingLeft: 0,
                        paddingTop: '2vh'
                    }}>
                        <MessageList
                            className='message-list'
                            lockable={true}
                            toBottomHeight={'100%'}
                            dataSource={[
                                {
                                    position: 'right',
                                    type: 'text',
                                    text: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit',
                                    date: new Date(),
                                }
                        ]} />

                    </div>
                    
                </div>
            </div>


            
        )
    }


}