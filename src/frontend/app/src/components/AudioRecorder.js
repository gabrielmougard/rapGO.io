import React, { Component } from 'react'
import { connect } from 'react-redux';
import { ReactMic } from 'react-mic'
import {Button} from 'baseui/button';
import TriangleRight from 'baseui/icon/triangle-right';
import Upload from 'baseui/icon/upload';
import Check from 'baseui/icon/check';

import { getRap } from '../actions';

require('./styles.scss')

class AudioRecorder extends Component {
  constructor(props) {
    super(props)
    this.state = {
      downloadLinkURL: null,
      isRecording: false,
      recordingStarted: false,
      recordingStopped: false,
      isProcessingBLOB: false,
      rawInputBLOB: null
    }
  }

  stopRecording= () => {
    this.setState({ isRecording: false })
  }

  onSave=(blobObject) => {
    console.log("on Save")
    this.setState({
      downloadLinkURL: blobObject.blobURL,
      rawInputBLOB: blobObject
    })
  }

  onStart=() => {
    console.log('You can tap into the onStart callback')
  }

  onStop= (blobObject) => {
    console.log("on Stop !")
    this.setState({ blobURL: blobObject.blobURL })
  }

  onData(recordedBlob){
    console.log('ONDATA CALL IS BEING CALLED! ', recordedBlob);
  }

  onBlock() {
    alert('ya blocked me!')
  }

  startRecording= () => {
    this.setState({
      isRecording: true,
      recordingInSession: true,
      recordingStarted: true,
      recordingStopped: false,
      isPaused: false
    })
  }

  stopRecording=() => {
    this.setState({
      isRecording: false,
      recordingInSession: false,
      recordingStarted: false,
      recordingStopped: true
    })
  }

  sendBLOB=() => {
    this.setState({
      isProcessingBLOB: true
    })
    this.props.getRap(this.state.rawInputBLOB);
  }

  render() {
    const {
      blobURL,
      downloadLinkURL,
      isRecording,
      recordingInSession,
      recordingStarted,
      recordingStopped
    } = this.state

    const recordBtn = recordingInSession ? true : false
    const downloadLink = recordingStopped ? "fa fa-download" : "fa disabled fa-download"
    
    return (
      <div>
        <div id="project-wrapper">
          <div id="project-container">
            <div id="overlay" />
            <div id="content">
              <h2>RapGO.io - generate a rap with your voice !</h2>
              <ReactMic
                className="oscilloscope"
                record={isRecording}
                backgroundColor="#333333"
                visualSetting="sinewave"
                audioBitsPerSecond={128000}
                onStop={this.onStop}
                onStart={this.onStart}
                onSave={this.onSave}
                onData={this.onData}
                onBlock={this.onBlock}
                onPause={this.onPause}
                strokeColor="#0096ef"
              />
              <div id="oscilloscope-scrim">
                {!recordingInSession && <div id="scrim" />}
              </div>
              <div id="controls">

                <Button isLoading={recordBtn} onClick={this.startRecording} startEnhancer={() => <TriangleRight size={30} />}>
                  Start Recording
                </Button>
                <span />
                <Button onClick={this.stopRecording} startEnhancer={() => <Check size={30} />}>
                  Stop Recording
                </Button>
                
                <div className="column download">
                  <Button onClick={this.sendBLOB} endEnhancer={() => <Upload size={30} />}>
                    Generate !
                </Button>
                </div>
              </div>
            </div>
            <div id="audio-playback-controls">
              <audio ref="audioSource" controls="controls" src={blobURL} controlsList="nodownload"/>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    isProcessingBLOB: state.isProcessingBLOB
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    getRap: (inputBLOB) => { dispatch(getRap(inputBLOB)) },
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(AudioRecorder);