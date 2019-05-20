module.exports = function (api) {
  
  const isTest = api.env('test');
  
  api.cache(true);

  const presets = ['@babel/react']

  if(isTest) {
    presets.push(["@babel/preset-env", {
      targets: {
        node: 'current',
      },
    }])
  } else {
    presets.push("@babel/preset-env")
  }

  const plugins = [];
  
  return {
    presets,
    plugins
  };
}
