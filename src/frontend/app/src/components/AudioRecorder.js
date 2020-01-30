import React, { Component } from 'react'
import { connect } from 'react-redux';
import { ReactMic } from 'react-mic'
import {Button} from 'baseui/button';
import {ProgressBar} from 'baseui/progress-bar';
import TriangleRight from 'baseui/icon/triangle-right';
import Upload from 'baseui/icon/upload';
import Check from 'baseui/icon/check';

import Websocket from 'react-websocket';

import { getRap } from '../actions';
import { downloadOutput } from '../actions';

import { NB_STEPS_HEARTBEATS } from '../CONSTANTS';
import { HEARTBEAT_TRIGGER_DOWNLOAD } from '../CONSTANTS';
import { HEARTBEAT_WS_ENDPOINT } from '../CONSTANTS';
import { NON_CHANGING_HEARTBEAT } from '../CONSTANTS';
import OutputLoader from './OutputLoader';
import OutputPlayer from './OutputPlayer';

require('./styleRecorder.scss')

class AudioRecorder extends Component {
  constructor(props) {
    super(props)
    this.state = {
      downloadLinkURL: null,
      isRecording: false,
      recordingStarted: false,
      recordingStopped: false,
      isProcessingBLOB: false,
      rawInputBLOB: null,
      heartbeatUUID: "",
      wsInterface: null,
      heartbeatMsg: "Waiting for input...",
      progressbarValue: 0, 
      deleteWsInterface: false,
      svgLoaderEnabled: false,
      outputLoaded: false
    }
    this.handleIncomingHeartbeat = this.handleIncomingHeartbeat.bind(this)

  }

  componentDidUpdate() {
    //wen the raw data is loaded from server delete the svg loader
    if (this.props.outputResponse && !this.state.outputLoaded) {
      console.log("output raw data detected !")
      this.setState({outputLoaded: true})
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

  handleIncomingHeartbeat(data) {
    let result = JSON.parse(data);
    console.log(result)
    let description = result.Description
    // not used for now
    //timestamp = result.Timestamp

    //update the progressBar
    if (description != NON_CHANGING_HEARTBEAT) {
      this.setState({progressbarValue: this.state.progressbarValue + 100/NB_STEPS_HEARTBEATS}) // the integer value
      this.setState({heartbeatMsg: description}) // the label below
    } else {
      this.setState({heartbeatMsg: description}) // just change the label

      //delete the wsInterface here
      this.setState({deleteWsInterface: true})
      //

    }

    
    if (description == HEARTBEAT_TRIGGER_DOWNLOAD) {
      this.setState({progressbarValue: 100})
      let uuid = this.props.heartbeatUUID.split("_")[1].split(".")[0];
      this.setState({svgLoaderEnabled: true})
      this.setState({deleteWsInterface: true})
      this.props.downloadOutput(uuid) // Amazing link !!! ==> https://apiko.com/blog/how-to-work-with-sound-java-script/
    }
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
    
    let output;

    //svg loader
    if (this.state.svgLoaderEnabled && !this.state.outputLoaded) {
      console.log("BERLUSCONI1")
      output = <div class="output-container">
                  <OutputLoader />
               </div>
    } else if (this.props.outputResponse && this.state.svgLoaderEnabled) {
      console.log("BERLUSCONI2")
      console.log(this.props.outputResponse)
      output = <div class="output-container">
                 <OutputPlayer response={this.props.outputResponse}/>
               </div>
    } else {
      output = <div class="output-container"></div>
    }

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
        {this.props.heartbeatUUID && !this.state.deleteWsInterface ? (
          <div class="progressbar-container">
          <ProgressBar
            value={this.state.progressbarValue}
            successValue={100}
            getProgressLabel={(currentValue, successValue) => this.state.heartbeatMsg}
            showLabel
          />
          <Websocket url={HEARTBEAT_WS_ENDPOINT+this.props.heartbeatUUID.split("_")[1].split(".")[0]} onMessage={this.handleIncomingHeartbeat} />
          </div>
        ) : (
          <div class="progressbar-container"></div>
        )}
        {output}
      </div>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    isProcessingBLOB: state.isProcessingBLOB,
    heartbeatUUID: state.heartbeatUUID,
    outputResponse: state.outputResponse
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    getRap: (inputBLOB) => { dispatch(getRap(inputBLOB)) },
    downloadOutput: (uuid) => {dispatch(downloadOutput(uuid))}
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(AudioRecorder);