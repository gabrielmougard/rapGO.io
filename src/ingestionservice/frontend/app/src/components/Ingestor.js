import React, { Component } from "react";
import {
  ThemeProvider,
  createTheme,
  lightThemePrimitives
} from "baseui";
import {
  HeaderNavigation,
  ALIGN,
  StyledNavigationList,
  StyledNavigationItem
} from "baseui/header-navigation";
import { Button, SIZE } from "baseui/button";
import { Display4 } from "baseui/typography";
import { Drawer } from 'baseui/drawer';
import Overflow from 'baseui/icon/overflow';
import {Textarea} from 'baseui/textarea';
import { Tag } from "baseui/tag";
import { FileUploader } from "baseui/file-uploader";

import Gauge from 'react-svg-gauge';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import Typography from '@material-ui/core/Typography';
import ListSubheader from '@material-ui/core/ListSubheader';
import grey from "@material-ui/core/colors/grey";
import { makeStyles } from '@material-ui/core/styles';

import MetadataLoader from './MetadataLoader';
import GaugeLoader from './GaugeLoader';

//import { connect } from 'react-redux';

import { INGESTOR_BUILD_VERSION } from "../CONSTANTS";

require('./styles.css')

class Ingestor extends Component {
  constructor(props) {
    super(props)
    this.state = {
      urlsDrawer : false,
      urlsDrawerTextAreaValue : "",
      urlsFileDescDrawer : false,
      rawfilesDrawer : false,
      genreDrawer : false,
      gaugeValue : 40
    }
  }

  closeURLSdrawer() {
    this.setState({urlsDrawer : false})
  }
  closeURLSfileDescDrawer() {
    this.setState({urlsFileDescDrawer : false})
  }
  closeRawFilesDrawer() {
    this.setState({rawfilesDrawer : false})
  }
  closeGenreDrawer() {
    this.setState({genreDrawer : false})
  }
  openURLSdrawer() {
    this.setState({urlsDrawer : true})
  }
  openURLSfileDescDrawer() {
    this.setState({urlsFileDescDrawer : true})
  }
  openRawFilesDrawer() {
    this.setState({rawfilesDrawer : true})
  }
  openGenreDrawer() {
    this.setState({genreDrawer : true})
  }
  getHexColor(value) {
    let string = value.toString(16);
    return (string.length === 1) ? '0' + string : string;
  }
  refreshMetadata() {

  }
  setTextAreaValueURLsDrawer(newValue) {
    this.setState({urlsDrawerTextAreaValue : newValue})
  }
  sendURLs() {

  }
  sendURLsFileDesc() {

  }

  render() {
    // const {} = this.state
    const tLabelStyle = {
      textAnchor: "middle",
      fill: "#999999",
      stroke: "none",
      fontStyle: "normal",
      fontVariant: "normal",
      fontWeight: 'normal',
      fontStretch: 'normal',
      lineHeight: 'normal',
      fillOpacity: 1,
      fontSize: 20
    }
    const labelStyle = {
      textAnchor: "middle",
      fill: "#010101",
      stroke: "none",
      fontStyle: "normal",
      fontVariant: "normal",
      fontWeight: 'normal',
      fontStretch: 'normal',
      lineHeight: 'normal',
      fillOpacity: 1
    }
    
    //for the gauge
    let r = Math.floor(this.state.gaugeValue * 2.55);
    let g = Math.floor(255 - (this.state.gaugeValue * 2.55));
    let b = 0;
    let colorHex = '#' + this.getHexColor(r) + this.getHexColor(g) + this.getHexColor(b);

    //for the cards
    const cardColor = 'rgb(255, 0, 255)'
    return (
    <div>

    <HeaderNavigation
      overrides={{
        Root: {
          style: ({ $theme }) => {
            return {};
          }
        }
      }}
    >

      <StyledNavigationList $align={ALIGN.left}>
        <StyledNavigationItem>
          <Display4>RapGO.io</Display4>
        </StyledNavigationItem>
      </StyledNavigationList>

      <StyledNavigationList $align={ALIGN.center}>
        <StyledNavigationItem>
            <Display4>Ingestion engine interface</Display4>
        </StyledNavigationItem>
      </StyledNavigationList>

      <StyledNavigationList $align={ALIGN.right}>
        <StyledNavigationItem>
          <Button>To GCP bucket</Button>
        </StyledNavigationItem>
      </StyledNavigationList>

    </HeaderNavigation>

    
    <div class="column-left">
      <div class="title-left-column">
        <Display4>Ingest data by :</Display4>
      </div>

      <div class="ingest-button-container">
        <Button
              onClick={() =>
                this.openURLSdrawer()
              }
              size={SIZE.large}
              overrides={{
                BaseButton: {
                  style: {
                    marginTop: '12px',
                    marginBottom: '12px',
                    marginLeft: '12px',
                    marginRight: '12px',
                  },
                },
              }}
        >
          URLs
        </Button>
        <Drawer
            onClose={() => this.closeURLSdrawer()}
            isOpen={this.state.urlsDrawer}
            anchor="left"
          >
            <Display4>Add the URLs to be ingested below :</Display4>
            <div class="textarea-urls-container">
              <Textarea
                value={this.state.urlsDrawerTextAreaValue}
                onChange={e => this.setTextAreaValueURLsDrawer(e.currentTarget.value)}
                placeholder="Enter one URL per line..."
                overrides={{
                  Input: {
                    style: {
                      maxHeight: '300px',
                      minHeight: '100px',
                      minWidth: '300px',
                      width: '100vw', // fill all available space up to parent max-width
                      resize: 'both',
                    },
                  },
                  InputContainer: {
                    style: {
                      maxWidth: '100%',
                      width: 'min-content',
                    },
                  },
                }}
              />
            </div>
            <Tag closeable={false} variant={"solid"} kind="accent">
            hint : https://freemusicarchive.org/EXAMPLE_SONG.mp3)
            </Tag>
            <div class="urls-drawer-button-container">
              <Button
                onClick={() =>
                  this.sendURLs()
                }
                size={SIZE.large}
                overrides={{
                  BaseButton: {
                    style: {
                      marginTop: '12px',
                      marginBottom: '12px',
                      marginLeft: '12px',
                      marginRight: '12px',
                    },
                  },
                }}
              >
                Ingest
              </Button>
            </div>
        </Drawer>
      </div>
      <div class="ingest-button-container">
      <Button
              onClick={() =>
                this.openURLSfileDescDrawer()
              }
              size={SIZE.large}
              overrides={{
                BaseButton: {
                  style: {
                    marginTop: '12px',
                    marginBottom: '12px',
                    marginLeft: '12px',
                    marginRight: '12px',
                  },
                },
              }}
        >
          URLs file descriptor
        </Button>
        <Drawer
            onClose={() => this.closeURLSfileDescDrawer()}
            isOpen={this.state.urlsFileDescDrawer}
            anchor="left"
          >
            <Display4>Add the URLs file description (text file containing the URLs, one per line) below :</Display4>
            <div class="filedroparea-urlsfiledesc-container">
              <FileUploader
                //onCancel={} //shoud put cancelUpload function in param
                onDrop={(acceptedFiles, rejectedFiles) => {
                  // handle file upload...
                  //startFakeProgress();
                }}
                // progressAmount is a number from 0 - 100 which indicates the percent of file transfer completed
                //progressAmount={progressAmount}
                //progressMessage={
                //  progressAmount
                //    ? `Uploading... ${progressAmount}% of 100%`
                //    : ''
                //}
              />
            </div>
            <div class="urlsfiledesc-drawer-button-container">
              <Button
                onClick={() =>
                  this.sendURLsFileDesc()
                }
                size={SIZE.large}
                overrides={{
                  BaseButton: {
                    style: {
                      marginTop: '12px',
                      marginBottom: '12px',
                      marginLeft: '12px',
                      marginRight: '12px',
                    },
                  },
                }}
              >
                Ingest
              </Button>
            </div>
        </Drawer>
      </div>
      <div class="ingest-button-container">
        <Button
              onClick={() =>
                this.openRawFilesDrawer()
              }
              size={SIZE.large}
              overrides={{
                BaseButton: {
                  style: {
                    marginTop: '12px',
                    marginBottom: '12px',
                    marginLeft: '12px',
                    marginRight: '12px',
                  },
                },
              }}
        >
          Raw files
        </Button>
        <Drawer
            onClose={() => this.closeRawFilesDrawer()}
            isOpen={this.state.rawfilesDrawer}
            anchor="left"
          >
            <Display4>Add the raw files (.mp3, .webm, ...) below :</Display4>
            <div class="filedroparea-rawfile-container">
              <FileUploader
                //onCancel={} //shoud put cancelUpload function in param
                onDrop={(acceptedFiles, rejectedFiles) => {
                  // handle file upload...
                  //startFakeProgress();
                }}
                // progressAmount is a number from 0 - 100 which indicates the percent of file transfer completed
                //progressAmount={progressAmount}
                //progressMessage={
                //  progressAmount
                //    ? `Uploading... ${progressAmount}% of 100%`
                //    : ''
                //}
              />
            </div>
            <div class="rawfile-drawer-button-container">
              <Button
                onClick={() =>
                  this.sendRawFiles()
                }
                size={SIZE.large}
                overrides={{
                  BaseButton: {
                    style: {
                      marginTop: '12px',
                      marginBottom: '12px',
                      marginLeft: '12px',
                      marginRight: '12px',
                    },
                  },
                }}
              >
                Ingest
              </Button>
            </div>
        </Drawer>
      </div>
      <div class="ingest-button-container">
        <Button
              onClick={() =>
                this.openGenreDrawer()
              }
              size={SIZE.large}
              overrides={{
                BaseButton: {
                  style: {
                    marginTop: '12px',
                    marginBottom: '12px',
                    marginLeft: '12px',
                    marginRight: '12px',
                  },
                },
              }}
        >
          Genre
        </Button>
        <Drawer
            onClose={() => this.closeGenreDrawer()}
            isOpen={this.state.genreDrawer}
            anchor="left"
          >
            Add by musical genre
        </Drawer>
      </div>
    </div>

    <div class="column-right">
        
        <div class="gauge-container">
          <div>
            <div class="metadata-header-title"><Display4>Bucket metadata</Display4></div>
            <div class="metadata-header-refresh">
              <Button
                onClick={() =>
                  this.refreshMetadata()
                }
                size={SIZE.large}
                overrides={{
                  BaseButton: {
                    style: {
                      marginTop: '12px',
                      marginBottom: '12px',
                      marginLeft: '12px',
                      marginRight: '12px',
                    },
                  },
                }}
                endEnhancer={() => <Overflow size={24} />}
              >
              Refresh metadata
              </Button>
            </div>
          </div>
          <Gauge value={this.state.gaugeValue} topLabelStyle={tLabelStyle} valueLabelStyle={labelStyle} width={400} height={320} color={colorHex} label="0.34Go/30Go consummed" valueFormatter={value => `${value}%`}/>
        </div>

        <GridList cellHeight={180}>
          <GridListTile key="Subheader" cols={2} style={{ height: 'auto' }}>
            <ListSubheader component="div">Objects statistics</ListSubheader>
          </GridListTile>
          <GridListTile>
              <Card>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    audio/mpeg
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"input_<UUID>.mp3"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    audio/mpeg
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"output_<UUID>.mp3"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"3.1Mo"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 562 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    audio/mpeg
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"beat_<UUID>.mp3"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    application/octet-stream
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"tempDist_<UUID>"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    application/octet-stream
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"tempInt_<UUID>"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    application/octet-stream
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"bpm_<UUID>"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    application/octet-stream
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"verseInterval_<UUID>"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>
          <GridListTile>
              <Card minWidth={100} style={{ cardColor }}>
                <CardContent>
                  <Typography fontSize={14} color="textSecondary" gutterBottom>
                    application/octet-stream
                  </Typography>
                  <Typography variant="h5" component="h2">
                    {"duration_<UUID>"}
                  </Typography>
                  <Typography marginBottom={12} color="textSecondary">
                    Avg. size : {"330Ko"}
                  </Typography>
                  <Typography variant="body2" component="p">
                    Occurence : 5642 units
                  </Typography>
                </CardContent>
              </Card>
          </GridListTile>           
        </GridList>
    </div>
    </div>
    )
  }
}

// const mapStateToProps = (state) => {
//   return {
//     isProcessingBLOB: state.isProcessingBLOB
//   }
// }

// const mapDispatchToProps = (dispatch) => {
//   return {
//     getRap: (inputBLOB) => { dispatch(getRap(inputBLOB)) },
//   }
// }

//export default connect(mapStateToProps, mapDispatchToProps)(Ingestor);

export default Ingestor;