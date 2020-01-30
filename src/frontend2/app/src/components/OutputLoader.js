import React, { Component } from "react";
import ContentLoader from "react-content-loader" 

class OutputLoader extends Component {
    render() {
        return (
            <ContentLoader 
                speed={1}
                width={400}
                height={475}
                viewBox="0 0 400 475"
                backgroundColor="#f3f3f3"
                foregroundColor="#ecebeb"
            >
                <rect x="43" y="245" rx="2" ry="2" width="140" height="18" /> 
                <rect x="207" y="245" rx="2" ry="2" width="140" height="18" /> 
                <rect x="0" y="60" rx="2" ry="2" width="400" height="172" />
            </ContentLoader>
        )
    }
}

export default OutputLoader