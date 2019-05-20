
export default class AjaxAssistant {

  constructor($) {
    this.$ = $
  }

  handleError(xhr, reject) {
    var text = ""
    try {
      const data = JSON.parse(xhr.responseText)
      text = data.errors
    } catch(e) {
      text = xhr.responseText
    }

    if(text === "") {
      if(xhr.status === 404) {
        text = "A resource could not be found."
      }
    }
    
    reject(text)
  }
  
  post(path, data) {
    return new Promise((resolve, reject) => {
      if(!data) data = {}
      this.$.ajax({
        method: 'POST',
        url: path,
        dataType: 'JSON',
        data: data,
        beforeSend: (xhr) => { xhr.setRequestHeader('CSRF-Token', this.$('meta[name=csrf-token]').attr('content')) },
        success: (data) => resolve(data),
        error: (xhr) => this.handleError(xhr, reject)
      })
    })
  }
  
  get(path, data) {
    return new Promise((resolve, reject) => {
      let ajaxParams = {
        url: path,
        dataType: 'JSON',
        success: (data) => resolve(data),
        error: (xhr) => this.handleError(xhr, reject)
      }

      if(data) {
        ajaxParams[data] = data
      }
      
      this.$.ajax(ajaxParams)
    })
  }
}
