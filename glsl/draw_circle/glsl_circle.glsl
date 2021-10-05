
#extension GL_OES_standard_derivatives : enable

#ifdef GL_ES
precision mediump float;
#endif

varying vec4 v_position;
varying vec4 v_normal;
varying vec2 v_texcoord;
varying vec4 v_color;

uniform mat4 u_projectionMatrix;
uniform mat4 u_modelViewMatrix;
uniform mat4 u_normalMatrix;
uniform vec2 u_resolution;
uniform float u_time;

#if defined(VERTEX)

// attribute vec4 a_position; // myfolder/myfile.obj
attribute vec4 a_position;
attribute vec4 a_normal;
attribute vec2 a_texcoord;
attribute vec4 a_color;

void main(void) {
	v_position = u_projectionMatrix * u_modelViewMatrix * a_position;
	v_normal = u_normalMatrix * a_normal;
	v_texcoord = a_texcoord;
	v_color = a_color;
	gl_Position = v_position;
}

#else // fragment shader

uniform vec2 u_mouse;
uniform vec2 u_pos;
// uniform sampler2D u_texture; // https://cdn.jsdelivr.net/gh/actarian/plausible-brdf-shader/textures/mars/4096x2048/diffuse.jpg?repeat=true
// uniform vec2 u_textureResolution;

#extension GL_OES_standard_derivatives : enable

#ifdef GL_ES
precision mediump float;
#endif

const float aWidth = 1.25;
const vec3 aColor = vec3(1,0,1);
const vec4 bkColor = vec4(1);

void main() {
    float dist = distance(u_mouse + vec2(100, -100), gl_FragCoord.xy);
    float r = distance(u_mouse, vec2(0, u_resolution.x))/5.0; // aRadius
    float delta = fwidth(dist);
    float inR = r - aWidth;
	float alphaIn = smoothstep(inR, inR - delta, dist);
    float alpha = smoothstep(r-delta, r, dist);
    vec4 circleColor = mix(vec4(aColor, 1), bkColor, alphaIn > 0.0 ? alphaIn : alpha);
	gl_FragColor = vec4(circleColor.rgb, circleColor.a);
}

#endif
