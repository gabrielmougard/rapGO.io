import React, { Component } from 'react';
import InputRange from 'react-input-range';
import 'react-input-range/lib/css/index.css'
import { loadRawData, getAudioContext } from '../lib/outputDecoder';

require('./stylePlayer.scss')
//https://apiko.com/blog/how-to-work-with-sound-java-script/
class OutputPlayer extends Component {
    constructor(props) {
      super(props)
      this.state = {
          volumeLevel: 100,
          progress: 0,
          playState: 'play',
          loading: false,
          player: null,
          audionState:{
              startedAt: null,
              pausedAt: null,
              isPause: true,
              duration: 0,
          } 
      }
      this.onProgressClick = this.onProgressClick.bind(this)
      this.onPlayBtnClick = this.onPlayBtnClick.bind(this)
      this.onStopBtnClick = this.onStopBtnClick.bind(this)
      this.onVolumeChange = this.onVolumeChange.bind(this)
    }

    componentDidMount() {
      setInterval(() => {
        if (this.state.audionState.startedAt && !this.state.audionState.isPause) {
          const playbackTime = (Date.now() - this.state.audionState.startedAt) / 1000;
          const rate = parseInt((playbackTime * 100) / this.state.audionState.duration, 10);
          rate <= 100 && this.setState({progress: rate})
        }
      }, 1000)
    }

    async onPlayBtnClick() {
        try {
            if(!this.state.player) {
                this.setState({loading: true})
                const frequencyC = document.querySelector('.frequency-bars');
                const sinewaveC = document.querySelector('.sinewave');
                const newPlayer = await loadRawData(this.props.response.data, 
                    { frequencyC, sinewaveC }, 
                    { fillStyle: 'rgb(250, 250, 250)', // background
                      strokeStyle: 'rgb(251, 89, 17)', // line color
                      lineWidth: 1,
                      fftSize: 16384 //delization of bars from 1024 to 32768
                    });
                
                this.setState({loading: false})
                this.setState({player: newPlayer})
                this.setState({audionState: {
                    startedAt: Date.now(),
                    isPause: false,
                    duration: newPlayer.duration,
                }});
                newPlayer.play(0)
                this.setState({playState: 'stop'})
                return 
            }
            this.setState({audionState: {
                startedAt: Date.now() - this.state.audionState.pausedAt,
                isPause: false,
            }});
            this.state.player.play(this.state.audionState.pausedAt / 1000);
            this.setState({playState: 'stop'})
            return
        } catch (e) {
            this.setState({loading: false})
            console.log(e);
        }
    }

    onStopBtnClick() {
     
      //https://github.com/VolodymyrTymets/sound-in-js/blob/master/client/src/components/Example4/Container.js
      this.setState({audionState: {
        pausedAt: Date.now() - this.state.audionState.startedAt,
        isPause: true,
      }});
      this.state.player && this.state.player.stop();
      this.setState({playState: 'play'})
    }

    onVolumeChange(max) {
      const value = max / 100;
      const level = value > 0.5 ? value * 4 : value * -4;
      this.state.player.setVolume(level || -1);
      this.setState({volumeLevel: max || 0})
    }

    onProgressClick(e) {
      const rate = (e.clientX * 100) / e.target.offsetWidth;
      const playbackTime = (this.state.audionState.duration * rate) / 100;

      this.state.player && this.state.player.stop();
      this.state.player && this.state.player.play(playbackTime);

      this.setState({progress: parseInt(rate, 10)});
      this.setState({audionState: {
        startedAt: Date.now() - playbackTime * 100,
      }});
    }

    render() {
        return (
            <div>
                <h4>Result : <small className="text-muted">(still alpha version so it's far from being perfect)</small></h4>

                <div className="bars-wrapper">
                  <canvas className="frequency-bars" width="512" height="100"></canvas>
                  <canvas className="sinewave" width="512" height="100"></canvas>
                </div>

                <div className="player mt-4">
                  <div className="progress player-progress mb-2" onClick={this.onProgressClick}>
                    <div
                      className="progress-bar bg-warning"
                      role="progressbar"
                      style={{width: `${this.state.progress}%`}}
                      aria-valuemax="100"
                      onClick={e => console.log(e)}
                    >
                    </div>
                  </div>
                  <div className="player-controls mt-2">

                    <div>{this.state.loading && <i className="fas fa-spinner fa-spin"></i>}</div>

                    <button
                      type="button"
                      className="btn btn-warning"
                      onClick={this.state.playState === 'play' ? this.onPlayBtnClick : this.onStopBtnClick}
                      disabled={this.state.loading}>
                      <i className={`fas fa-${this.state.playState}`}></i>
                    </button>

                    <div className="player-volume-control">
                      <i onClick={() => this.onVolumeChange(0)} className="fas fa-volume-down"></i>
                      <div className="range-select">
                        <InputRange
                          maxValue={100}
                          minValue={0}
                          value={this.state.volumeLevel}
                          onChange={this.onVolumeChange}
                        />
                      </div>
                      <i onClick={() => this.onVolumeChange(100)}  className="fas fa-volume-up"></i>
                    </div>
                  </div>
                </div>
            </div>
        )
    }
}

export default OutputPlayer;