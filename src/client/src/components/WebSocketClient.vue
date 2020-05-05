<template>
    <div></div>
</template>

<script>
export default {
  name: 'WebsocketClient',
  data () {
    return {
      title: 'Websocket Echo Client',
      message: '',
      logs: [],
      status: 'disconnected'
    }
  },
  created () {
    this.connect()
  },
  methods: {
    connect () {
      this.socket = new WebSocket('wss://echo.websocket.org')
      this.socket.onopen = () => {
        this.status = 'connected'
        console.log('WebSocket connected to:', this.socket.url)
        this.logs.push({event: 'WebSocket Connected', data: this.socket.url})
        this.socket.onmessage = ({data}) => {
          this.logs.push({event: 'Recieved message', data})
          console.log('Received:', data)
        }
      }
    },
    disconnect () {
      this.socket.close()
      this.status = 'disconnected'
      this.logs = []
      console.log('WebSocket disconnected')
    },
    sendMessage () {
      // Send message to Websocket echo service
      this.socket.send(this.message)
      this.logs.push({ event: 'Sent message', data: this.message })
      // POST message user enetered to our server
      this.axios.post(process.env.BASE_URL + 'log-collector', {message: this.message})
        .catch(function (error) {
          console.log(error)
      })
      // Log to console and clear input field
      console.log('Sent:', this.message)
      this.message = ''
    }
  }
}
</script>