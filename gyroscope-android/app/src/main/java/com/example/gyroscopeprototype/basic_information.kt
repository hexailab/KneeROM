package com.example.gyroscopeprototype

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup

/**
 * A simple [Fragment] subclass.
 * Use the [basic_information.newInstance] factory method to
 * create an instance of this fragment.
 */
class basic_information : Fragment() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        arguments?.let {}
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_basic_information, container, false)
    }

    companion object {
        @JvmStatic
        fun newInstance() =
            basic_information().apply {
                arguments = Bundle().apply {}
            }
    }
}