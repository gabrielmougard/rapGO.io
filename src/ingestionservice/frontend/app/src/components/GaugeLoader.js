import React, { Component } from "react";
import ContentLoader from "react-content-loader" 

class GaugeLoader extends Component {
    render() {
        return (
            <ContentLoader 
                height={320}
                width={400}
                speed={2}
                backgroundColor="#a6a4a4"
                primaryColor="#a6a4a4"
                secondaryColor="#f5f5f5"
            >
                <rect x="0" y="0" rx="3" ry="3" width="400" height="320" />
            </ContentLoader>
        )
    }
}

export default GaugeLoader